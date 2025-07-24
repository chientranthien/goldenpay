import React from 'react';
import { Contact, User } from '../api/chatService';

interface ContactsListProps {
    contacts: Contact[];
    selectedContact: Contact | null;
    onContactSelect: (contact: Contact) => void;
    loading: boolean;
}

const ContactsList: React.FC<ContactsListProps> = ({
    contacts,
    selectedContact,
    onContactSelect,
    loading
}) => {
    const formatLastSeen = (timestamp: number) => {
        const date = new Date(timestamp);
        const now = new Date();
        const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60));

        if (diffInMinutes < 1) {
            return 'Just now';
        } else if (diffInMinutes < 60) {
            return `${diffInMinutes}m ago`;
        } else if (diffInMinutes < 1440) {
            return `${Math.floor(diffInMinutes / 60)}h ago`;
        } else {
            return date.toLocaleDateString();
        }
    };

    if (loading && contacts.length === 0) {
        return (
            <div className="contacts-loading">
                <div className="loading-spinner"></div>
                <p>Loading contacts...</p>
            </div>
        );
    }

    if (contacts.length === 0) {
        return (
            <div className="contacts-empty">
                <p>No contacts found</p>
                <small>Add contacts to start chatting</small>
            </div>
        );
    }

    return (
        <div className="contacts-list">
            {contacts.map((contact) => {
                if (!contact) return null;

                const isSelected = selectedContact?.id === contact.id

                return (
                    <div
                        key={contact.id}
                        className={`contact-item ${isSelected ? 'selected' : ''}`}
                        onClick={() => onContactSelect(contact)}
                    >
                        <div className="contact-avatar">
                            <div className="avatar-circle">
                                {contact.name.charAt(0).toUpperCase()}
                            </div>
                            <div className="status-indicator online"></div>
                        </div>

                        <div className="contact-info">
                            <div className="contact-name">{contact.name}</div>
                            <div className="contact-email">{contact.email}</div>
                        </div>

                        <div className="contact-actions">
                            <div className="unread-badge">0</div>
                        </div>
                    </div>
                );
            })}

            {loading && (
                <div className="contacts-loading-more">
                    <div className="loading-spinner small"></div>
                </div>
            )}
        </div>
    );
};

export default ContactsList; 