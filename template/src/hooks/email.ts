import { useCallback } from "react";
import { emails } from "../types/email";

const useEmail = () => {
  const getEmails = useCallback(async () => {
    const response = await fetch("/emails");
    const data: emails = await response.json();
    return data;
  }, []);

  return { getEmails };
};

export default useEmail;
