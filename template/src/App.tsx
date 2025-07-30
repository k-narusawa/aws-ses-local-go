import { useEffect, useState } from "react";
import { emails } from "./types/email";
import useEmail from "./hooks/email";
import EmailTable from "./components/EmailTable";

function App() {
  const { getEmails, deleteEmail, deleteEmails } = useEmail();
  const [emails, setEmails] = useState<emails | null>(null);
  const [page, setPage] = useState<number>(1);
  const [searchAddress, setSearchAddress] = useState<string>("");

  useEffect(() => {
    const fetchEmails = async () => {
      const resp = await getEmails(page, 20, searchAddress);
      setEmails(resp);
    };
    fetchEmails();
    const intervalId = setInterval(fetchEmails, 10000);

    return () => clearInterval(intervalId);
  }, [getEmails, page, searchAddress, deleteEmail, deleteEmails]);

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold text-gray-900">
            AWS SES Local メール管理
          </h1>
        </div>
      </header>
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {!emails ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
          </div>
        ) : (
          <div className="bg-white shadow rounded-lg">
            <EmailTable
              emails={emails}
              setEmails={setEmails}
              setPage={setPage}
              setSearchAddress={setSearchAddress}
              deleteEmail={deleteEmail}
              deleteEmails={deleteEmails}
            />
          </div>
        )}
      </main>
    </div>
  );
}

export default App;
