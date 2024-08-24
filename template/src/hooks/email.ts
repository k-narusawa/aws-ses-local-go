import { useCallback } from "react";
import { emails } from "../types/email";
import axios, { AxiosResponse } from "axios";

const useEmail = () => {
  const getEmails = useCallback(async (page: number, limit: number) => {
    const response = await axios
      .get(`/emails`, {
        params: {
          page: page,
          limit: limit,
        },
      })
      .then((res: AxiosResponse<emails>) => res.data);

    return response;
  }, []);

  return { getEmails };
};

export default useEmail;
