

export type Notification = {
  id?: String // id would be used to remove timeout, may only need to be local type
  type: "error" | "success" | "info"
  message: String
  detail?: String
}

export type Notifications = Array<Notification>
