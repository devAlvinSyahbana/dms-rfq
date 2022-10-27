import { createContext, useEffect, useRef, useState } from "react";
import "./App.css";
import { gql } from "graphql-request";
import { request, GraphQLClient } from "graphql-request";
import { Formik, FieldArray, Field } from "formik";
const query = (email: string, password: string) => gql`
  mutation {
    login(input: { email: "${email}", password: "${password}" }) {
      token
    }
  }
`;

const queryHome = gql`
  query {
    RFQList {
      id
      CompanyName
      QuotationNo
    }
  }
`;

const queryNew = (val) => gql`
    mutation {
      createRFQ(
        input: {
          Tax: ${val.Tax}
          SNK: "${val.SNK}"
          CompanyName:"${val.CompanyName}"
          CompanyAddress:"${val.CompanyAddress}"
          CompanyWebsite:"${val.CompanyWebsite}"
          SentToName:"${val.SentToName}"
          SentToPhone:"${val.SentToPhone}"
          SentToAddress:"${val.SentToAddress}"
          MadeForName:"${val.MadeForName}"
          MadeForPhone:"${val.MadeForPhone}"
          MadeForAddress:"${val.MadeForAddress}"
          Disc:${val.Disc}
          Interest:${val.Interest}
          QuotationDate:"${val.QuotationDate}"
          QuotationNo:"${val.QuotationNo}"
          QuotationExpires:"${val.QuotationExpires}"
          Items:${JSON.stringify(val.Items).replace(/"([^"]+)":/g, "$1:")}
        }
      ) {
        id
      }
    }
`;

const AppContext = createContext({});
function App() {
  const [loggedIn, setLoggedIn] = useState({ status: false });
  const [loginFailed, setLoginFailed] = useState(false);
  const [page, setPage] = useState<string>("home");
  const [gqlClient, setGqlClient] = useState<any>(
    new GraphQLClient(import.meta.env.VITE_GQL_HOST, {
      headers: { Authorization: "Bearer " + localStorage.getItem("token") },
    })
  );
  const [homeState, setHomeState] = useState<any[]>();

  const inputRef = useRef<any>();
  const passwordRef = useRef<any>();

  const handleLoggedIn = async () => {
    await setGqlClient(
      new GraphQLClient(import.meta.env.VITE_GQL_HOST, {
        headers: { Authorization: "Bearer " + localStorage.getItem("token") },
      })
    );
    setLoggedIn((prev) => ({ ...prev, status: true }));
  };
  useEffect(() => {
    if (localStorage && localStorage.getItem("token") !== "") {
      handleLoggedIn();
    }
  }, []);

  const getListTable = async () => {
    const homeQuery = await gqlClient.request({
      url: import.meta.env.VITE_GQL_HOST,
      document: queryHome,
      requestHeaders: {
        Authorization: "Bearer " + localStorage.getItem("token"),
      },
    });
    setHomeState(homeQuery.RFQList);
  };
  const insertData = async (v) => {
    await gqlClient.request({
      url: import.meta.env.VITE_GQL_HOST,
      document: queryNew(v),
      requestHeaders: {
        Authorization: "Bearer " + localStorage.getItem("token"),
      },
    });
    setPage("home");
  };

  useEffect(() => {
    if (page === "home" && gqlClient) {
      getListTable();
    }
  }, [page, gqlClient]);

  const handleLogin = async (e: any) => {
    e.preventDefault();
    const email = inputRef.current.value;
    const password = passwordRef.current.value;
    await request(import.meta.env.VITE_GQL_HOST, query(email, password))
      .then(async (data) => {
        setLoginFailed(false);
        await setGqlClient(
          new GraphQLClient(import.meta.env.VITE_GQL_HOST, {
            headers: { Authorization: "Bearer " + data.login.token },
          })
        );
        localStorage.setItem("token", data.login.token);
        setLoggedIn((prev) => ({ ...prev, status: true }));
      })
      .catch((e) => {
        setLoginFailed(true);
      });
  };

  const handleLogout = async (e: any) => {
    e.preventDefault();
    localStorage.removeItem("token");
    setLoggedIn((prev) => ({ ...prev, status: false }));
  };

  const Home = () => {
    return (
      <>
        <div className="overflow-x-auto relative m-20">
          <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
            <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
              <tr>
                <th scope="col" className="py-3 px-6">
                  Company Name
                </th>
                <th scope="col" className="py-3 px-6">
                  Quotation No
                </th>
                <th scope="col" className="py-3 px-6">
                  Action
                </th>
              </tr>
            </thead>
            <tbody>
              {homeState &&
                homeState.length > 0 &&
                homeState.map((item: any) => {
                  return (
                    <tr className="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                      <th
                        scope="row"
                        className="py-4 px-6 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                      >
                        {item.CompanyName}
                      </th>
                      <td className="py-4 px-6"> {item.QuotationNo}</td>
                      <td className="py-4 px-6">
                        <a
                          className="bg-indigo-600 p-2 rounded text-white"
                          href={`http://localhost:4080/download/${item.id}`}
                          target="_blank"
                        >
                          Detail
                        </a>
                      </td>
                    </tr>
                  );
                })}
            </tbody>
          </table>
        </div>
      </>
    );
  };

  const Template = ({ children }: any) => {
    return (
      <div className="bg-gray-800 w-screen min-h-screen">
        <nav className="bg-white border-gray-200 px-2 sm:px-4 py-2.5  dark:bg-gray-900">
          <div className="container flex flex-wrap justify-between items-center mx-auto">
            <div className="w-full md:block md:w-auto" id="navbar-default">
              <ul className="flex flex-col p-4 mt-4 bg-gray-50  border border-gray-100 md:flex-row md:space-x-8 md:mt-0 md:text-sm md:font-medium md:border-0 md:bg-white dark:bg-gray-800 md:dark:bg-gray-900 dark:border-gray-700">
                <li>
                  <a
                    href="#"
                    onClick={(e) => setPage("home")}
                    className="block py-2 pr-4 pl-3 text-white bg-indigo-700 rounded md:bg-transparent md:text-indigo-700 md:p-0 dark:text-white"
                    aria-current="page"
                  >
                    Home
                  </a>
                </li>
                {loggedIn && loggedIn.status && (
                  <>
                    <li>
                      <a
                        href="#"
                        onClick={(e) => setPage("new")}
                        className="block py-2 pr-4 pl-3 text-gray-700 rounded hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 md:p-0 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent"
                      >
                        New RFQ
                      </a>
                    </li>
                  </>
                )}
                <li>
                  <a
                    onClick={handleLogout}
                    className="block py-2 pr-4 pl-3 text-gray-700 rounded hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 md:p-0 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent"
                  >
                    Logout
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </nav>
        {children}
      </div>
    );
  };

  const WText = ({ children }) => {
    return <span className="text-white">{children}</span>;
  };
  const NewPage = () => {
    return (
      <div className="m-6">
        <Formik
          initialValues={{
            Tax: 10,
            SNK: ["SNK 1", "SNK 2"],
            CompanyName: "Tripatra",
            CompanyAddress: "Bintaro",
            CompanyWebsite: "Tripatra.com",
            SentToName: "Hendra",
            SentToPhone: "021456789",
            SentToAddress: "Bandung",
            MadeForName: "Ziat",
            MadeForPhone: "0123456789",
            MadeForAddress: "123456789",
            Disc: 10,
            Interest: 10,
            QuotationDate: "2022-10-27",
            QuotationNo: "123456789",
            QuotationExpires: "2022-10-30",
            Items: [
              { Nama: "A", Harga: 1000, Qty: 10 },
              { Nama: "B", Harga: 1000, Qty: 10 },
            ],
          }}
          onSubmit={(values, { setSubmitting }) => {
            insertData(values);
            setSubmitting(false);
          }}
        >
          {({
            values,
            errors,
            touched,
            handleChange,
            handleBlur,
            handleSubmit,
            isSubmitting,
          }) => (
            <form onSubmit={handleSubmit} className="flex flex-col">
              <WText>Tax</WText>
              <input
                className="rounded"
                name="Tax"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.Tax}
              />
              <WText>CompanyName</WText>
              <input
                className="rounded"
                name="CompanyName"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.CompanyName}
              />
              <WText>CompanyAddress</WText>
              <input
                className="rounded"
                name="CompanyAddress"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.CompanyAddress}
              />
              <WText>CompanyWebsite</WText>
              <input
                className="rounded"
                name="CompanyWebsite"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.CompanyWebsite}
              />
              <WText>SentToName</WText>
              <input
                className="rounded"
                name="SentToName"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.SentToName}
              />
              <WText>SentToAddress</WText>
              <input
                className="rounded"
                name="SentToAddress"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.SentToAddress}
              />
              <WText>MadeForName</WText>
              <input
                className="rounded"
                name="MadeForName"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.MadeForName}
              />
              <WText>MadeForPhone</WText>
              <input
                className="rounded"
                name="MadeForPhone"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.MadeForPhone}
              />
              <WText>MadeForAddress</WText>
              <input
                className="rounded"
                name="MadeForAddress"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.MadeForAddress}
              />
              <WText>Disc</WText>
              <input
                className="rounded"
                name="Disc"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.Disc}
              />
              <WText>Interest</WText>
              <input
                className="rounded"
                name="Interest"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.Interest}
              />
              <WText>QuotationDate </WText>
              <input
                className="rounded"
                name="QuotationDate"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.QuotationDate}
              />
              <WText>QuotationNo </WText>
              <input
                className="rounded"
                name="QuotationNo"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.QuotationNo}
              />
              <WText>QuotationExpires </WText>
              <input
                className="rounded"
                name="QuotationExpires"
                onChange={handleChange}
                onBlur={handleBlur}
                value={values.QuotationExpires}
              />
              <WText>SNK </WText>

              <FieldArray
                name="SNK"
                render={(arrayHelpers) => (
                  <div>
                    {values.SNK && values.SNK.length > 0 ? (
                      values.SNK.map((friend, index) => (
                        <div key={index}>
                          <Field name={`SNK.${index}`} />

                          <button
                            type="button"
                            onClick={() => arrayHelpers.remove(index)} // remove a friend from the list
                            className="text-white"
                          >
                            -
                          </button>

                          <button
                            type="button"
                            onClick={() => arrayHelpers.insert(index, "")} // insert an empty string at a position
                            className="text-white"
                          >
                            +
                          </button>
                        </div>
                      ))
                    ) : (
                      <button
                        type="button"
                        onClick={() => arrayHelpers.push("")}
                      >
                        Add SNK
                      </button>
                    )}
                  </div>
                )}
              />
              <WText>Items</WText>

              <FieldArray
                name="Items"
                render={(arrayHelpers) => (
                  <div>
                    {values.Items.map((friend, index) => (
                      <div key={index} className="flex">
                        {/** both these conventions do the same */}
                        <WText>Nama</WText>
                        <Field name={`Items[${index}].Nama`} />
                        <WText>Harga</WText>
                        <Field type="number" name={`Items.${index}.Harga`} />
                        <WText>Qty</WText>
                        <Field type="number" name={`Items.${index}.Qty`} />
                        <button
                          type="button"
                          onClick={() => arrayHelpers.remove(index)}
                        >
                          -
                        </button>
                      </div>
                    ))}

                    <button
                      type="button"
                      onClick={() =>
                        arrayHelpers.push({ Nama: "", Harga: 0, Qty: 0 })
                      }
                    >
                      +
                    </button>
                  </div>
                )}
              />
              <button
                type="submit"
                className="text-white p-2 mt-2 bg-indigo-700 rounded"
                disabled={isSubmitting}
              >
                Submit
              </button>
            </form>
          )}
        </Formik>
      </div>
    );
  };

  if (loggedIn.status) {
    if (gqlClient !== undefined) {
      return (
        <Template>
          {page === "home" && <Home />}
          {page === "new" && <NewPage />}
        </Template>
      );
    }
  } else {
    return (
      <div className="w-screen h-screen bg-gray-700 flex items-center">
        <div className="w-1/2 bg-slate-100 m-auto rounded flex flex-col">
          <span className="font-bold text-center align-middle m-auto grow">
            Login
          </span>

          <div className="flex flex-col px-72">
            <div className="flex flex-col mt-5">
              <label htmlFor="">Email</label>
              <input type="text" className="border " ref={inputRef} />
            </div>
            <div className="flex flex-col mt-5">
              <label htmlFor="">Password</label>
              <input type="password" className="border " ref={passwordRef} />
            </div>

            {loginFailed && (
              <span className="font-bold text-center align-middle m-auto grow bg-red-400 rounded px-2 mt-5">
                Invalid email / password
              </span>
            )}

            <div className="flex flex-col mt-5 mb-14">
              <button
                className="bg-gray-700 rounded text-white"
                onClick={handleLogin}
              >
                Login
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default App;
