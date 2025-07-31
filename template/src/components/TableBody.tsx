import { useState } from "react";
import { emails } from "../types/email";
import React from "react";

type Props = {
  emails: emails;
  deleteEmail: (message_id: string) => void;
};

export const TableBody: React.FC<Props> = ({ emails, deleteEmail }) => {
  const [openRows, setOpenRows] = useState<{ [key: string]: boolean }>({});

  const onDelete = (message_id: string) => {
    deleteEmail(message_id);
    emails.items = emails.items.filter(
      (email) => email.message_id !== message_id
    );
  };

  const toggleRow = (message_id: string) => {
    setOpenRows((prev) => ({
      ...prev,
      [message_id]: !prev[message_id],
    }));
  };

  return (
    <tbody className="divide-y divide-gray-200 bg-white">
      {emails?.items.map((email) => (
        <React.Fragment key={email.message_id}>
          <tr
            onClick={() => toggleRow(email.message_id)}
            className="hover:bg-gray-50 transition-colors duration-150 cursor-pointer"
          >
            <td className="py-3 ps-4">
              <div className="flex items-center h-5">
                <input
                  type="checkbox"
                  className="w-4 h-4 border-gray-300 rounded text-blue-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200"
                  aria-label={`Select ${email.subject}`}
                  onClick={(e) => e.stopPropagation()}
                />
              </div>
            </td>
            <td className="px-6 py-4 text-sm font-medium text-gray-900 max-w-xs truncate">
              {email.subject || "(No subject)"}
            </td>
            <td className="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">
              {email.from}
            </td>
            <td className="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">
              {email.destination.to}
            </td>
            <td className="px-6 py-4 text-end text-sm font-medium">
              <button
                type="button"
                className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-red-600 hover:text-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transition-colors duration-200"
                onClick={(e) => {
                  e.stopPropagation();
                  onDelete(email.message_id);
                }}
                aria-label={`Delete ${email.subject}`}
              >
                <svg
                  className="h-4 w-4 mr-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
                Delete
              </button>
            </td>
          </tr>
          {openRows[email.message_id] && (
            <tr className="bg-gray-50">
              <td colSpan={5} className="px-6 py-4">
                <div className="space-y-4">
                  <div>
                    <h4 className="text-sm font-medium text-gray-900 mb-2">
                      Text Content:
                    </h4>
                    <div className="bg-white rounded-md shadow-sm">
                      <pre className="p-4 text-sm text-gray-700 whitespace-pre-wrap">
                        {email.body.text || "(No text content)"}
                      </pre>
                    </div>
                  </div>
                  <div>
                    <h4 className="text-sm font-medium text-gray-900 mb-2">
                      HTML Content:
                    </h4>
                    <div
                      className="bg-white rounded-md shadow-sm p-4 prose prose-sm max-w-none"
                      dangerouslySetInnerHTML={{
                        __html: email.body.html || "(No HTML content)",
                      }}
                    />
                  </div>
                </div>
              </td>
            </tr>
          )}
        </React.Fragment>
      ))}
    </tbody>
  );
};
