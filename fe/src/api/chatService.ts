import { config, ToJSON } from './common'

// API Response interfaces
interface ApiResponse {
    code: Code;
    data?: any;
    message?: string;
}

interface Code {
    id: number;
    msg: string;
}

// Chat types
interface Contact {
    id: number;
    contactUserId: number;
    name: string;
    email: string;
}

interface User {
    id: number;
    email: string;
    name: string;
    status: number;
    ctime: number;
    mtime: number;
}


interface Channel {
    id: number;
    name: string;
    description: string;
    type: string;
    creatorId: number;
    status: number;
    ctime: number;
    mtime: number;
}

interface Message {
    id: number;
    channelId: number;
    userId: number;
    content: string;
    messageType: string;
    status: number;
    ctime: number;
    mtime: number;
}

interface GetContactsRequest {
    pagination?: {
        val: number;
        limit: number;
        hasMore: boolean;
    };
}

interface SendMessageRequest {
    channelId: number;
    content: string;
    messageType?: string;
}

interface GetMessagesRequest {
    channelId: number;
    pagination?: {
        val: number;
        limit: number;
        hasMore: boolean;
    };
}

interface CreateDirectMessageRequest {
    toUserId: number;
}

interface IChatService {
    GetContacts: () => Promise<{ contacts: Contact[], code: Code }>;
    CreateDirectMessage: (toUserId: number) => Promise<{ channelId: number, code: Code }>;
    SendMessage: (channelId: number, content: string) => Promise<{ messageId: number, code: Code }>;
    GetMessages: (channelId: number, pagination?: any) => Promise<{ messages: Message[], code: Code }>;
    GetDirectMessages: (pagination?: any) => Promise<{ channels: Channel[], code: Code }>;
}

const chatService: IChatService = {
    GetContacts: async function (): Promise<{ contacts: Contact[], code: Code }> {
        const api = config.host + "users/contacts/_query"
        const body: GetContactsRequest = {
            pagination: {
                val: 0,
                hasMore: true,
                limit: 50,
            }
        }

        let contacts: Contact[] = []
        let code: Code = { id: 0, msg: "" }

        await fetch(api, {
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json"
            },
            body: ToJSON(body)
        }).then(resp => {
            return resp.json()
        }).then((json: ApiResponse) => {
            code = json.code
            if (json.data && json.data.contacts) {
                contacts = json.data.contacts
            }
        }).catch(() => {
            console.log("failed to get contacts")
        })

        return { contacts, code }
    },

    CreateDirectMessage: async function (toUserId: number): Promise<{ channelId: number, code: Code }> {
        const api = config.host + "api/v1/chat/direct-messages"
        const body = { toUserId }

        let channelId: number = 0
        let code: Code = { id: 0, msg: "" }

        await fetch(api, {
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json"
            },
            body: ToJSON(body)
        }).then(resp => {
            return resp.json()
        }).then((json: ApiResponse) => {
            code = json.code
            if (json.data && json.data.channelId) {
                channelId = json.data.channelId
            }
        }).catch(() => {
            console.log("failed to create direct message")
        })

        return { channelId, code }
    },

    SendMessage: async function (channelId: number, content: string): Promise<{ messageId: number, code: Code }> {
        const api = config.host + "api/v1/chat/messages"
        const body: SendMessageRequest = {
            channelId,
            content,
            messageType: "text"
        }

        let messageId: number = 0
        let code: Code = { id: 0, msg: "" }

        await fetch(api, {
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json"
            },
            body: ToJSON(body)
        }).then(resp => {
            return resp.json()
        }).then((json: ApiResponse) => {
            code = json.code
            if (json.data && json.data.messageId) {
                messageId = json.data.messageId
            }
        }).catch(() => {
            console.log("failed to send message")
        })

        return { messageId, code }
    },

    GetMessages: async function (channelId: number, pagination?: any): Promise<{ messages: Message[], code: Code }> {
        const api = config.host + `api/v1/chat/channels/${channelId}/messages`
        const body = {
            pagination: pagination || { limit: 50 }
        }

        let messages: Message[] = []
        let code: Code = { id: 0, msg: "" }

        await fetch(api, {
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json"
            },
            body: ToJSON(body)
        }).then(resp => {
            return resp.json()
        }).then((json: ApiResponse) => {
            code = json.code
            if (json.data && json.data.messages) {
                messages = json.data.messages
            }
        }).catch(() => {
            console.log("failed to get messages")
        })

        return { messages, code }
    },

    GetDirectMessages: async function (): Promise<{ channels: Channel[], code: Code }> {
        const api = config.host + "api/v1/chat/direct-messages"
        let channels: Channel[] = []
        let code: Code = { id: 0, msg: "" }

        await fetch(api, {
            method: 'GET',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json"
            }
        }).then(resp => {
            return resp.json()
        }).then((json: ApiResponse) => {
            code = json.code
            if (json.data && json.data.channels) {
                channels = json.data.channels
            }
        }).catch(() => {
            console.log("failed to get direct messages")
        })

        return { channels, code }
    }
}

export default function ChatService(): IChatService {
    return chatService
}

export type { Contact, User, Channel, Message } 