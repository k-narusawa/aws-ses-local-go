import { useEffect, useState } from "react";
import { emails } from "./types/email";
import useEmail from "./hooks/email";
import EmailTable from "./components/EmailTable";

function App() {
  const { getEmails } = useEmail();
  const [emails, setEmails] = useState<emails | null>(null);
  const [page, setPage] = useState<number>(1);

  useEffect(() => {
    (async () => {
      const resp = await getEmails(page, 10);
      setEmails(resp);
    })();
  }, [getEmails, page]);

  if (!emails) {
    return <div>Loading...</div>;
  }

  return (
    <div className="p-4">
      <EmailTable emails={emails} />
    </div>
  );
}

export default App;
