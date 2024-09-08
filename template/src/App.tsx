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
    (async () => {
      const resp = await getEmails(page, 10, searchAddress);
      setEmails(resp);
    })();
  }, [getEmails, page, searchAddress, deleteEmail, deleteEmails]);

  if (!emails) {
    return <div>Loading...</div>;
  }

  return (
    <div className="p-4">
      <EmailTable
        emails={emails}
        setEmails={setEmails}
        setPage={setPage}
        setSearchAddress={setSearchAddress}
        deleteEmail={deleteEmail}
        deleteEmails={deleteEmails}
      />
    </div>
  );
}

export default App;
