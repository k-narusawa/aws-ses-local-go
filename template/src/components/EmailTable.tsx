import { emails } from "../types/email";
import { PaginationNav } from "./PaginationNav";
import { SearchInput } from "./SearchInput";
import { TableBody } from "./TableBody";
import { TableHeader } from "./TableHeader";

type Props = {
  emails: emails;
  setPage: (page: number) => void;
  setSearchAddress: (address: string) => void;
};

const EmailTable: React.FC<Props> = ({ emails, setPage, setSearchAddress }) => {
  return (
    <>
      <div className="flex flex-col">
        <div className="-m-1.5 overflow-x-auto">
          <div className="p-1.5 min-w-full inline-block align-middle">
            <div className="border rounded-lg divide-y divide-gray-200">
              <div className="py-3 px-4">
                <div className="relative max-w-xs">
                  <SearchInput setSearchAddress={setSearchAddress} />
                </div>
              </div>

              <div className="overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <TableHeader />
                  <TableBody emails={emails} />
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