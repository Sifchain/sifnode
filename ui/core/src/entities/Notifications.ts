export type Notification = {
  id?: string; // id would be used to remove timeout, may only need to be local type
  type: "error" | "success" | "info";
  message: string;
  detail?: string;
};

export type Notifications = Array<Notification>;
