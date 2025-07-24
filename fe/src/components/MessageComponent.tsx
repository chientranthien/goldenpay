import React from 'react';
import { Message } from '../api/chatService';

interface MessageComponentProps {
    message: Message;
    isOwn: boolean;
    showTime: boolean;
    contactName: string;
}

const MessageComponent: React.FC<MessageComponentProps> = ({
    message,
    isOwn,
    showTime,
    contactName
}) => {
    const formatTime = (timestamp: number) => {
        return new Date(timestamp).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const formatMessageContent = (content: string) => {
        // Basic text formatting - can be extended for rich text, links, etc.
        return content.split('\n').map((line, index) => (
            <React.Fragment key={index}>
                {line}
                {index < content.split('\n').length - 1 && <br />}
            </React.Fragment>
        ));
    };

    const getMessageStatus = (): 'pending' | 'delivered' | 'read' | null => {
        // For now, just show delivered for sent messages
        if (isOwn) {
            return message.status === 1 ? 'delivered' : 'pending';
        }
        return null;
    };

    const messageStatusIcon = () => {
        const status = getMessageStatus();
        if (!status) return null;

        switch (status) {
            case 'pending':
                return <span className="message-status pending" title="Sending">⏳</span>;
            case 'delivered':
                return <span className="message-status delivered" title="Delivered">✓</span>;
            case 'read':
                return <span className="message-status read" title="Read">✓✓</span>;
            default:
                return null;
        }
    };

    return (
        <div className={`message ${isOwn ? 'own' : 'other'}`}>
            {!isOwn && showTime && (
                <div className="message-sender">
                    <div className="sender-avatar">
                        {contactName.charAt(0).toUpperCase()}
                    </div>
                    <span className="sender-name">{contactName}</span>
                </div>
            )}

            <div className="message-content">
                <div className={`message-bubble ${message.messageType || 'text'}`}>
                    {message.messageType === 'text' ? (
                        <div className="message-text">
                            {formatMessageContent(message.content)}
                        </div>
                    ) : message.messageType === 'file' ? (
                        <div className="message-file">
                            <div className="file-icon">📎</div>
                            <div className="file-info">
                                <div className="file-name">{message.content}</div>
                                <div className="file-size">Click to download</div>
                            </div>
                        </div>
                    ) : message.messageType === 'system' ? (
                        <div className="message-system">
                            <em>{message.content}</em>
                        </div>
                    ) : (
                        <div className="message-text">
                            {formatMessageContent(message.content)}
                        </div>
                    )}

                    <div className="message-meta">
                        <span className="message-time">
                            {formatTime(message.ctime)}
                        </span>
                        {messageStatusIcon()}
                    </div>
                </div>

                {message.mtime !== message.ctime && (
                    <div className="message-edited">
                        <small>edited</small>
                    </div>
                )}
            </div>
        </div>
    );
};

export default MessageComponent; 