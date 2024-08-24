import { useEffect, useState } from "react";
import { emails } from "./types/email";
import useEmail from "./hooks/email";
import EmailTable from "./components/EmailTable";

function App() {
  const { getEmails } = useEmail();
  const [emails, setEmails] = useState<emails | null>(null);

  useEffect(() => {
    (async () => {
      const resp = await getEmails();
      setEmails(resp);
    })();
  }, [getEmails]);

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
