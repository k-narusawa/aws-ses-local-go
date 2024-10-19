import { useCallback } from "react";
import { emails } from "../types/email";
import axios, { AxiosResponse } from "axios";

const baseURL = import.meta.env.VITE_API_URL ?? undefined;

const useEmail = () => {
  const getEmails = useCallback(
    async (page: number, limit: number, toAddress: string) => {
      const response = await axios
        .get(`/emails`, {
          baseURL: baseURL,
          params: {
            page: page,
            limit: limit,
            to_address: toAddress,
          },
        })
        .then((res: AxiosResponse<emails>) => res.data);

      return response;
    },
    []
  );

  const deleteEmail = useCallback(async (id: string) => {
    const response = await axios
      .delete(`/emails/${id}}`, {
        baseURL: baseURL,
      })
      .then((res) => res.data);

    return response;
  }, []);

  const deleteEmails = useCallback(async () => {
    const response = await axios
      .delete(`/emails`, {
        baseURL: baseURL,
      })
      .then((res) => res.data);

    return response;
  }, []);

  return { getEmails, deleteEmail, deleteEmails };
};

export default useEmail;
