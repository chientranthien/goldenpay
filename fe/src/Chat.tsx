import React, { useState, useEffect } from 'react';
import ChatService, { Channel, Message, Contact } from './api/chatService';
import { useRedirectToLoginIfNotAuthenticated, GetUserIdFromCookie } from './common';
import ContactsList from './components/ContactsList';
import ChatBox from './components/ChatBox';
import Nav from './Nav';
import './Chat.css';

interface ChatState {
    contacts: Contact[];
    selectedContact: Contact | null;
    selectedChannel: Channel | null;
    messages: Message[];
    loading: boolean;
}

export default function Chat() {
    const { isAuthenticated } = useRedirectToLoginIfNotAuthenticated();
    const [chatState, setChatState] = useState<ChatState>({
        contacts: [],
        selectedContact: null,
        selectedChannel: null,
        messages: [],
        loading: false
    });

    const userId = GetUserIdFromCookie();

    useEffect(() => {
        if (isAuthenticated && userId) {
            loadContacts();
        }
    }, [isAuthenticated, userId]);

    const loadContacts = async () => {
        if (!userId) return;

        setChatState(prev => ({ ...prev, loading: true }));

        try {
            const { contacts, code } = await ChatService().GetContacts();
            if (code.id === 0) {
                setChatState(prev => ({ ...prev, contacts: contacts }));
            }
        } catch (error) {
            console.error('Failed to load contacts:', error);
        } finally {
            setChatState(prev => ({ ...prev, loading: false }));
        }
    };

    const handleContactSelect = async (contact: Contact) => {
        setChatState(prev => ({
            ...prev,
            selectedContact: contact,
            loading: true,
            messages: []
        }));

        try {
            // Create or get direct message channel
            const { channelId, code } = await ChatService().CreateDirectMessage(contact.contactUserId);
            if (code.id === 0) {
                const channel = {
                    id: channelId,
                    name: `dm_${userId}_${contact.contactUserId}`,
                    description: '',
                    type: 'direct',
                    creatorId: userId || 0,
                    status: 1,
                    ctime: Date.now(),
                    mtime: Date.now()
                };
                setChatState(prev => ({ ...prev, selectedChannel: channel }));

                // Load messages for this channel
                await loadMessages(channelId);
            }
        } catch (error) {
            console.error('Failed to create direct message:', error);
        } finally {
            setChatState(prev => ({ ...prev, loading: false }));
        }
    };

    const loadMessages = async (channelId: number) => {
        try {
            const { messages, code } = await ChatService().GetMessages(channelId);
            if (code.id === 0) {
                setChatState(prev => ({ ...prev, messages: messages.reverse() })); // Reverse to show oldest first
            }
        } catch (error) {
            console.error('Failed to load messages:', error);
        }
    };

    const handleSendMessage = async (content: string) => {
        if (!chatState.selectedChannel || !content.trim()) return;

        try {
            const { messageId, code } = await ChatService().SendMessage(chatState.selectedChannel.id, content);
            if (code.id === 0) {
                // Add the new message to the messages list
                const newMessage: Message = {
                    id: messageId,
                    channelId: chatState.selectedChannel.id,
                    userId: userId || 0,
                    content,
                    messageType: 'text',
                    status: 1,
                    ctime: Date.now(),
                    mtime: Date.now()
                };
                setChatState(prev => ({
                    ...prev,
                    messages: [...prev.messages, newMessage]
                }));
            }
        } catch (error) {
            console.error('Failed to send message:', error);
        }
    };

    if (isAuthenticated === null) {
        return <div className="chat-loading">Loading...</div>;
    }

    return (
        <>
            <Nav />
            <div className="chat-container">
                <div className="chat-sidebar">
                    <div className="chat-header">
                        <h3>Contacts</h3>
                    </div>
                    <ContactsList
                        contacts={chatState.contacts}
                        selectedContact={chatState.selectedContact}
                        onContactSelect={handleContactSelect}
                        loading={chatState.loading}
                    />
                </div>

                <div className="chat-main">
                    {chatState.selectedContact ? (
                        <ChatBox
                            contact={chatState.selectedContact}
                            messages={chatState.messages}
                            onSendMessage={handleSendMessage}
                            currentUserId={userId || 0}
                            loading={chatState.loading}
                        />
                    ) : (
                        <div className="chat-welcome">
                            <h2>Welcome to GoldenPay Chat</h2>
                            <p>Select a contact to start chatting</p>
                        </div>
                    )}
                </div>
            </div>
        </>
    );
} 