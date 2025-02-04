import { KeyIcon } from "@heroicons/react/24/outline";
import Nav from "components/Nav";
import { useState, useEffect } from "react";
import { format } from "date-fns";
//TODO: make all 3 transactions recent down
//TODO: fix date field for banking.  last digit of time gets cut off
export default function Administration() {
  //Shop data
  const [orders, setOrders] = useState();
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch(process.env.NEXT_PUBLIC_clientordersapi, {
        method: "GET",
      });
      const data = await response.json();
      setOrders(data);
    };
    fetchData().catch(console.error);
  }, []);
  //Insurance Data
  const [quotes, setQuotes] = useState();
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch(process.env.NEXT_PUBLIC_clientquotesapi, {
        method: "GET",
      });
      const data = await response.json();
      setQuotes(data.data.data);
    };
    fetchData().catch(console.error);
  }, []);
  //Banking Data
  const [transactions, setTransactions] = useState();
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch(process.env.NEXT_PUBLIC_clientmainapi + "/transactions", {
        method: "GET",
      });
      const data = await response.json();
      console.log(data);
      setTransactions(data);
    };
    fetchData().catch(console.error);
  }, []);
  const timezoneOffset = new Date().getTimezoneOffset();

  console.log(timezoneOffset);
  return (
    <div className="bg-gray-100">
      <Nav />
      <div className="flex p-4 py-2">
        <KeyIcon className="w-8 " />
        <h1 className=" text-4xl font-bold ml-2">Administration</h1>
      </div>
      <div className="flex gap-4 justify-around p-2">
        <div className="bg-white rounded-lg  shadow p-2 w-full">
          <span className="p-1 ">Shop Orders</span>

          <div className="overflow-x-auto relative mt-2 shadow-sm sm:rounded-lg">
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400 ">
              <thead className="text-xs text-gray-700 uppercase bg-gray-200">
                <tr>
                  <th scope="col" className="p-1">
                    Status
                  </th>
                  <th scope="col" className="p-1">
                    Order
                  </th>
                  <th scope="col" className="p-1">
                    Customer
                  </th>
                  <th scope="col" className="p-1">
                    Total Cost
                  </th>
                </tr>
              </thead>
              <tbody>
                {orders?.map((order) => (
                  <tr key={order.id} className="border-b bg-gray-100 ">
                    <th scope="row" className="p-1 font-medium text-gray-900 whitespace-nowrap ">
                      {order.status}
                    </th>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">#{order.id}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">{order.name}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">${order.cartTotal}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
        <div className="bg-white rounded-lg  shadow p-2 w-full">
          <span className="p-1 ">Banking Transactions</span>

          <div className="overflow-x-auto relative shadow-sm mt-2 sm:rounded-lg">
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400 ">
              <thead className="text-xs text-gray-700 uppercase bg-gray-200">
                <tr>
                  <th scope="col" className="p-1">
                    Customer
                  </th>
                  <th scope="col" className="p-1">
                    Account
                  </th>
                  <th scope="col" className="p-1">
                    Vendor
                  </th>
                  <th scope="col" className="p-1">
                    Amount
                  </th>
                  <th scope="col" className="p-1">
                    Date
                  </th>
                </tr>
              </thead>
              <tbody>
                {transactions?.map((transaction) => (
                  <tr key={transaction.id} className="border-b bg-gray-100 ">
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap ">{transaction.userId}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">{transaction.accountName}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">{transaction.vendor}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">${transaction.amount}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap ">
                      {format(new Date(transaction.timestamp), "M/d/yy H:m")}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
        <div className="bg-white rounded-lg  shadow  p-2 w-full">
          <span className="p-1 ">Quote Management</span>
          <div className="overflow-x-auto relative mt-2  sm:rounded-lg">
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400 ">
              <thead className="text-xs text-gray-700 uppercase bg-gray-200">
                <tr>
                  <th scope="col" className="p-1">
                    Status
                  </th>
                  <th scope="col" className="p-1">
                    Customer
                  </th>
                  <th scope="col" className="p-1">
                    Quote
                  </th>
                  <th scope="col" className="p-1">
                    Type
                  </th>
                </tr>
              </thead>
              <tbody>
                {quotes?.map((quote) => (
                  <tr key={quote.id} className="border-b bg-gray-100 ">
                    <td className="py-1 font-medium text-gray-900 whitespace-nowrap">{quote.status}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">{quote.name}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap ">{quote.id.slice(0, 7)}</td>
                    <td className="p-1 font-medium text-gray-900 whitespace-nowrap">Home/Auto</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
