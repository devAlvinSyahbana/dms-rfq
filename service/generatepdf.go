package service

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/devAlvinSyahbana/golang-rfq/graph/model"
	"github.com/gorilla/mux"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/lib/pq"
)

type ResultCalc struct {
	discount   int
	interest   int
	pajak      int
	totalharga int
}

func GeneratePDFMux(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := &model.Rfq{}

		vars := mux.Vars(r)

		err := db.QueryRow(`SELECT * FROM rfq.header WHERE "ID" = ($1)`, vars["id"]).Scan(
			&response.CompanyName,
			&response.CompanyAddress,
			&response.CompanyWebsite,
			&response.QuotationDate,
			&response.QuotationNo,
			&response.QuotationExpires,
			&response.MadeForName,
			&response.MadeForAddress,
			&response.MadeForPhone,
			&response.SentToName,
			&response.SentToAddress,
			&response.SentToPhone,
			&response.Disc,
			&response.Tax,
			&response.Interest,
			pq.Array(&response.Snk),
			&response.ID,
		)

		rows, err := db.Query(`SELECT * FROM rfq.items WHERE "HeaderID" = ($1)`, vars["id"])
		if err != nil {
			panic(err)
		}
		responseArray := []*model.Item{}
		for rows.Next() {
			responseItem := &model.Item{}
			rows.Scan(&responseItem.HeaderID,
				&responseItem.Nama,
				&responseItem.Harga,
				&responseItem.Qty)
			responseArray = append(responseArray, responseItem)
		}
		response.Items = responseArray
		if err != nil {
			panic(err)
		}
		GeneratePDF(response)
		f, err := os.Open("generated-" + response.QuotationNo + ".pdf")
		if f != nil {
			defer f.Close()
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contentDisposition := "attachment; filename=generated-" + response.QuotationNo + ".pdf"
		w.Header().Set("Content-Disposition", contentDisposition)
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func GeneratePDF(data *model.Rfq) {
	begin := time.Now()

	// darkGrayColor := getDarkGrayColor()getTotal
	// grayColor := getGrayColor()
	// whiteColor := color.NewWhite()
	blueColor := getBlueColor()
	// redColor := getRedColor()
	header := getHeader()
	contents := getContents(data.Items)
	total := getTotal(data.Items, data.Tax, data.Disc, data.Interest)
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.RegisterHeader(func() {
		m.Row(30, func() {
			m.Col(3, func() {
				m.Text(fmt.Sprintf("Nama Perusahaan : %s", data.CompanyName), props.Text{
					Size:        8,
					Align:       consts.Left,
					Extrapolate: false,
					Color:       blueColor,
				})
				m.Text(fmt.Sprintf("Alamat : %s", data.CompanyAddress), props.Text{
					Top:   10,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
				m.Text(fmt.Sprintf("Website: %s", data.CompanyWebsite), props.Text{
					Top:   15,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
			})

			m.Col(7, func() {

				m.Text("Tanggal", props.Text{
					Top:         8,
					Size:        8,
					Style:       consts.BoldItalic,
					Align:       consts.Right,
					Extrapolate: false,
					Color:       blueColor,
				})
				m.Text("No Quotation", props.Text{
					Top:   12,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: blueColor,
				})
				m.Text("Berlaku sampai", props.Text{
					Top:   16,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: blueColor,
				})
			})
			m.Col(1, func() {
				m.Text("Quotation", props.Text{
					Size:        14,
					Extrapolate: false,
				})
			})
			m.Col(3, func() {

				m.Text(fmt.Sprintf("%s", data.QuotationDate), props.Text{
					Top:         8,
					Size:        8,
					Style:       consts.BoldItalic,
					Extrapolate: false,
					Color:       blueColor,
				})
				m.Text(fmt.Sprintf("%s", data.QuotationNo), props.Text{
					Top:   12,
					Style: consts.BoldItalic,
					Size:  8,
					Color: blueColor,
				})
				m.Text(fmt.Sprintf("%s", data.QuotationExpires), props.Text{
					Top:   16,
					Style: consts.BoldItalic,
					Size:  8,
					Color: blueColor,
				})
			})
		})
		m.Row(30, func() {
			m.Col(3, func() {
				m.Text("Dibuat Untuk", props.Text{
					Size:        16,
					Align:       consts.Left,
					Extrapolate: false,
					Color:       blueColor,
				})
				m.Text(fmt.Sprintf("Nama Klien : %s", data.MadeForName), props.Text{
					Top:   10,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
				m.Text(fmt.Sprintf("Alamat: %s", data.MadeForAddress), props.Text{
					Top:   15,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
				m.Text(fmt.Sprintf("No Telepon : %s", data.MadeForPhone), props.Text{
					Top:   20,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
			})
			m.Col(7, func() {

				m.Text("Dikirimkan Ke", props.Text{
					Size:        16,
					Align:       consts.Right,
					Extrapolate: false,
					Color:       blueColor,
				})
				m.Text(fmt.Sprintf("Nama Klien : %s", data.SentToName), props.Text{
					Top:   10,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: blueColor,
				})
				m.Text(fmt.Sprintf("Alamat :%s", data.SentToAddress), props.Text{
					Top:   15,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: blueColor,
				})
				m.Text(fmt.Sprintf("No Telepon: %s", data.SentToPhone), props.Text{
					Top:   20,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: blueColor,
				})
			})
		})
	})

	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:              consts.Center,
		HeaderContentSpace: 1,
		Line:               true,
	})

	m.Row(20, func() {
		m.Col(3, func() {
			m.Text("Syarat dan Ketentuan", props.Text{
				Top:   5,
				Style: consts.BoldItalic,
				Size:  8,
				Align: consts.Left,
				Color: blueColor,
			})
			for i, j := range data.Snk {
				m.Text(j, props.Text{
					Top:   10 + float64(i*5),
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
			}
		})
		m.Col(20, func() {
			m.ColSpace(15)
			m.TableList([]string{"", "", ""}, total, props.TableList{
				HeaderProp: props.TableListContent{
					Size:      8,
					GridSizes: []uint{8, 4, 4},
				},
				ContentProp: props.TableListContent{
					Size:      8,
					GridSizes: []uint{10, 2, 2},
				},
				Align:              consts.Right,
				HeaderContentSpace: 1,
				Line:               false,
			})
		})
	})

	m.RegisterFooter(func() {
		m.Row(20, func() {
			m.Col(12, func() {
				m.Text("Terima kasih atas kepercayaan Anda!", props.Text{
					Top:   5,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Center,
					Color: blueColor,
				})
				m.Text("Jika ada pertanyaan lebih lanjut, silahkan hubungi xxx atau email xxx", props.Text{
					Top:   10,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Center,
					Color: blueColor,
				})
			})
		})

	})
	err := m.OutputFileAndClose("generated-" + data.QuotationNo + ".pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))
}

func getHeader() []string {
	return []string{"Barang/Jasa", "Harga per Unit", "Jumlah", "Total Harga"}
}

func getQuotation() [][]string {
	return [][]string{
		{"Tanggal", "2022-10-21"},
		{"No Quotation", "12312133"},
		{"Berlaku Sampai", "14 Hari"},
	}
}

func getResult(subTotal int, disc int, tax int, intr int) *ResultCalc {
	discount := subTotal - (subTotal * disc / 100)
	interest := subTotal - (subTotal * intr / 100)
	pajak := subTotal - (subTotal * tax / 100)
	totalharga := (subTotal - discount) + interest + pajak
	return &ResultCalc{
		discount:   discount,
		interest:   interest,
		pajak:      pajak,
		totalharga: totalharga,
	}

}
func getTotal(data []*model.Item, tax int, disc int, intr int) [][]string {

	subTotal := 0
	for _, row := range data {
		subTotal = subTotal + row.Qty*row.Harga
	}
	a := getResult(subTotal, disc, tax, intr)
	return [][]string{
		{"SubTotal", strconv.Itoa(subTotal)},
		{"Diskon", strconv.Itoa(a.discount)},
		{"Bunga Pajak", strconv.Itoa(a.interest)},
		{"Pajak", strconv.Itoa(a.pajak)},
		{"Total Harga", strconv.Itoa(a.totalharga)},
	}
}
func getContents(data []*model.Item) [][]string {
	items := [][]string{}
	for _, row := range data {
		item := []string{row.Nama, fmt.Sprintf("%d", row.Harga), fmt.Sprintf("%d", row.Qty), fmt.Sprintf("%d", row.Qty*row.Harga)}
		items = append(items, item)
	}
	return items
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlueColor() color.Color {
	return color.Color{
		Red:   10,
		Green: 10,
		Blue:  150,
	}
}

func getRedColor() color.Color {
	return color.Color{
		Red:   150,
		Green: 10,
		Blue:  10,
	}
}
