import { useState } from "react";
import { emails } from "../types/email";
import React from "react";

type Props = {
  emails: emails;
};

export const TableBody: React.FC<Props> = ({ emails }) => {
  const [openRows, setOpenRows] = useState<{ [key: string]: boolean }>({});

  const toggleRow = (message_id: string) => {
    setOpenRows((prev) => ({
      ...prev,
      [message_id]: !prev[message_id],
    }));
  };

  return (
    <>
      <tbody className="divide-y divide-gray-200">
        {emails && (
          <>
            {emails.items.map((email) => (
              <React.Fragment key={email.message_id}>
                <tr
                  onClick={() => toggleRow(email.message_id)}
                  className="hover:bg-gray-100"
                >
                  <td className="py-3 ps-4">
                    <div className="flex items-center h-5">
                      <input
                        id="hs-table-pagination-checkbox-2"
                        type="checkbox"
                        className="border-gray-200 rounded text-blue-600 focus:ring-blue-500"
                      />
                      <label className="sr-only">Checkbox</label>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-800">
                    {email.subject}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-800">
                    {email.from}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-800">
                    {email.to}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-end text-sm font-medium">
                    <button
                      type="button"
                      className="inline-flex items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent text-blue-600 hover:text-blue-800 focus:outline-none focus:text-blue-800 disabled:opacity-50 disabled:pointer-events-none"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
                {openRows[email.message_id] && (
                  <tr>
                    <td colSpan={6} className="px-6 py-4">
                      <div className="text-sm text-gray-800">
                        <div className="mb-2">
                          <strong>Text:</strong>
                          <textarea
                            className="w-full mt-2 p-2 border border-gray-300 rounded"
                            rows={4}
                            readOnly
                            value={email.text}
                          />
                        </div>
                        <div className="mb-2">
                          <strong>HTML:</strong>
                          <div
                            className="w-full mt-2 p-2 border border-gray-300 rounded"
                            dangerouslySetInnerHTML={{ __html: email.html }}
                          />
                        </div>
                      </div>
                    </td>
                  </tr>
                )}
              </React.Fragment>
            ))}
          </>
        )}
      </tbody>
    </>
  );
};
