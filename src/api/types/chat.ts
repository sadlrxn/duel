export interface Message {
  id: number;
  author: ChatUser;
  message: string;
  time: number;
  replyTo?: number;
  sponsors?: number[];
  deleted?: boolean;
}

export interface ChatUser {
  id: number;
  avatar: string;
  name: string;
  role?: string;
}
