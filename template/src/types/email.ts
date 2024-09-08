export type emails = {
  page: number;
  limit: number;
  size: number;
  total_size: number;
  total_page: number;
  items: email[];
};

export type email = {
  message_id: string;
  from: string;
  to: string;
  cc: string;
  bcc: string;
  subject: string;
  text: string;
  html: string;
  list_unsubscribe_post: string;
  list_unsubscribe_url: string;
  created_at: string;
};

export const emptyEmails = (): emails => {
  return {
    page: 0,
    limit: 0,
    size: 0,
    total_size: 0,
    total_page: 0,
    items: [],
  };
};
