import React, { useState, useRef, useEffect } from 'react';
import { Contact, Message } from '../api/chatService';
import MessageComponent from './MessageComponent';

interface ChatBoxProps {
    contact: Contact;
    messages: Message[];
    onSendMessage: (content: string) => void;
    currentUserId: number;
    loading: boolean;
}

const ChatBox: React.FC<ChatBoxProps> = ({
    contact,
    messages,
    onSendMessage,
    currentUserId,
    loading
}) => {
    const [newMessage, setNewMessage] = useState('');
    const [isSending, setIsSending] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLTextAreaElement>(null);

    // Auto scroll to bottom when new messages arrive
    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    // Focus input when contact changes
    useEffect(() => {
        if (inputRef.current) {
            inputRef.current.focus();
        }
    }, [contact.id]);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    const handleSendMessage = async () => {
        if (!newMessage.trim() || isSending) return;

        setIsSending(true);
        try {
            await onSendMessage(newMessage.trim());
            setNewMessage('');
        } catch (error) {
            console.error('Failed to send message:', error);
        } finally {
            setIsSending(false);
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSendMessage();
        }
    };

    const formatTime = (timestamp: number) => {
        return new Date(timestamp).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const groupMessagesByDate = (messages: Message[]) => {
        const groups: { [key: string]: Message[] } = {};

        messages.forEach(message => {
            const date = new Date(message.ctime).toDateString();
            if (!groups[date]) {
                groups[date] = [];
            }
            groups[date].push(message);
        });

        return groups;
    };

    const messageGroups = groupMessagesByDate(messages);
    const today = new Date().toDateString();
    const yesterday = new Date(Date.now() - 24 * 60 * 60 * 1000).toDateString();

    const formatDateGroup = (dateString: string) => {
        if (dateString === today) return 'Today';
        if (dateString === yesterday) return 'Yesterday';
        return new Date(dateString).toLocaleDateString();
    };

    return (
        <div className="chat-box">
            <div className="chat-box-header">
                <div className="contact-info">
                    <div className="contact-avatar">
                        <div className="avatar-circle">
                            {contact.name.charAt(0).toUpperCase()}
                        </div>
                        <div className="status-indicator online"></div>
                    </div>
                    <div className="contact-details">
                        <h3>{contact.name}</h3>
                        <p>{contact.email}</p>
                    </div>
                </div>
                <div className="chat-actions">
                    <button className="action-btn" title="Call">📞</button>
                    <button className="action-btn" title="Video Call">📹</button>
                    <button className="action-btn" title="More">⋯</button>
                </div>
            </div>

            <div className="chat-messages">
                {loading && messages.length === 0 ? (
                    <div className="messages-loading">
                        <div className="loading-spinner"></div>
                        <p>Loading messages...</p>
                    </div>
                ) : Object.keys(messageGroups).length === 0 ? (
                    <div className="messages-empty">
                        <div className="empty-state">
                            <h4>Start a conversation</h4>
                            <p>Send a message to {contact.name} to begin chatting</p>
                        </div>
                    </div>
                ) : (
                    <>
                        {Object.entries(messageGroups).map(([date, groupMessages]) => (
                            <div key={date} className="message-group">
                                <div className="date-divider">
                                    <span>{formatDateGroup(date)}</span>
                                </div>
                                {groupMessages.map((message, index) => {
                                    const isOwn = message.userId === currentUserId;
                                    const showTime = index === 0 ||
                                        groupMessages[index - 1].userId !== message.userId ||
                                        (message.ctime - groupMessages[index - 1].ctime) > 300000; // 5 minutes

                                    return (
                                        <MessageComponent
                                            key={message.id}
                                            message={message}
                                            isOwn={isOwn}
                                            showTime={showTime}
                                            contactName={contact.name}
                                        />
                                    );
                                })}
                            </div>
                        ))}
                    </>
                )}
                <div ref={messagesEndRef} />
            </div>

            <div className="chat-input">
                <div className="input-container">
                    <button className="attachment-btn" title="Attach file">📎</button>
                    <textarea
                        ref={inputRef}
                        value={newMessage}
                        onChange={(e) => setNewMessage(e.target.value)}
                        onKeyPress={handleKeyPress}
                        placeholder={`Message ${contact.name}...`}
                        className="message-input"
                        rows={1}
                        disabled={isSending}
                    />
                    <button
                        className={`send-btn ${newMessage.trim() ? 'active' : ''}`}
                        onClick={handleSendMessage}
                        disabled={!newMessage.trim() || isSending}
                        title="Send message"
                    >
                        {isSending ? '⏳' : '➤'}
                    </button>
                </div>
                <div className="input-actions">
                    <button className="emoji-btn" title="Add emoji">😊</button>
                    <span className="typing-indicator">
                        {isSending && 'Sending...'}
                    </span>
                </div>
            </div>
        </div>
    );
};

export default ChatBox; 