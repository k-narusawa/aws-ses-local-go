import { useEffect, useState } from "react";
import { emails } from "./types/email";
import "./App.css";
import useEmail from "./hooks/email";

function App() {
  const { getEmails } = useEmail();
  const [emails, setEmails] = useState<emails | null>(null);

  useEffect(() => {
    (async () => {
      const resp = await getEmails();
      setEmails(resp);
    })();
  }, [getEmails]);

  return (
    <>
      <h1>Emails</h1>
      {emails && (
        <ul>
          {emails.items.map((email) => (
            <li key={email.message_id}>
              <h2>{email.subject}</h2>
              <p>{email.text}</p>
            </li>
          ))}
        </ul>
      )}
    </>
  );
}

export default App;
