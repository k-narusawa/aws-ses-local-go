import { emails, emptyEmails } from "../types/email";
import { PaginationNav } from "./PaginationNav";
import { SearchInput } from "./SearchInput";
import { TableBody } from "./TableBody";
import { TableHeader } from "./TableHeader";

type Props = {
  emails: emails;
  setEmails: (emails: emails) => void;
  setPage: (page: number) => void;
  setSearchAddress: (address: string) => void;
  deleteEmail: (id: string) => void;
  deleteEmails: () => void;
};

const EmailTable: React.FC<Props> = ({
  emails,
  setEmails,
  setPage,
  setSearchAddress,
  deleteEmail,
  deleteEmails,
}) => {
  const onDeleteEmails = () => {
    deleteEmails();
    const newEmails = emptyEmails();
    setEmails(newEmails);
  };

  return (
    <div className="flex flex-col">
      <div className="p-6 border-b border-gray-200">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <SearchInput setSearchAddress={setSearchAddress} />
          <button
            type="button"
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-200"
            onClick={onDeleteEmails}
          >
            <svg
              className="h-5 w-5 mr-2"
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
            Delete All
          </button>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <TableHeader />
          <TableBody emails={emails} deleteEmail={deleteEmail} />
        </table>
      </div>

      <div className="p-6 border-t border-gray-200">
        <PaginationNav
          page={emails.page}
          totalPage={emails.total_page + 1}
          setPage={setPage}
        />
      </div>
    </div>
  );
};

export default EmailTable;
