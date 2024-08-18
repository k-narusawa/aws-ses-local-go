export type emails = {
  page: number;
  limit: number;
  size: number;
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
