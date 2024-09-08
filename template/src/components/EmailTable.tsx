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
    <>
      <div className="flex flex-col">
        <div className="-m-1.5 overflow-x-auto">
          <div className="p-1.5 min-w-full inline-block align-middle">
            <div className="border rounded-lg divide-y divide-gray-200">
              <div className="py-3 px-4">
                <div className="flex justify-between space-x-2">
                  <SearchInput setSearchAddress={setSearchAddress} />
                  <button
                    type="button"
                    className="focus:outline-none text-white bg-red-700 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-900"
                    onClick={onDeleteEmails}
                  >
                    Delete All
                  </button>
                </div>
              </div>

              <div className="overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <TableHeader />
                  <TableBody emails={emails} deleteEmail={deleteEmail} />
                </table>
              </div>

              <div className="py-1 px-4">
                <PaginationNav
                  page={emails.page}
                  totalPage={emails.total_page + 1}
                  setPage={setPage}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default EmailTable;
