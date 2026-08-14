import { useState } from "react";

import Chat from "@/components/chat/Chat";

import { useChatSession } from "@/context/ChatSessionContext";
import { useSelectedKnowledge } from "@/context/SelectedKnowledgeContext";

import type { ChatMode } from "@/types/chat-session";


export default function ChatPanel() {

    const {
        currentSessionId,
        currentSession,
        messages,
        isSending,
        isLoadingMessages,
        sendMessage,
    } = useChatSession();

    const { selected } = useSelectedKnowledge();

    const selectedNote = selected?.type === "note" ? selected.item : null;

    const [newChatMode, setNewChatMode] = useState<ChatMode>("knowledge");

    const isExistingSession = currentSessionId !== null;
    const chatMode: ChatMode = currentSession?.mode ?? newChatMode;

    async function handleSend(message: string) {

        if (isExistingSession) {
            await sendMessage(message);
            return;
        }

        await sendMessage(message, {
            mode: chatMode,
            noteId: chatMode === "note" ? selectedNote?.id : undefined,
            title:
                chatMode === "note" && selectedNote
                    ? selectedNote.title
                    : undefined,
        });
    }

    return (

        <section className="dashboard-panel chat-panel">

            <header className="panel-header">
                <div>
                    <h2>MINDK CHAT</h2>

                    <span>
                        {chatMode === "knowledge"
                            ? "Ask questions about your knowledge base."
                            : "Ask questions about the selected note."}
                    </span>
                </div>

                <button>...</button>
            </header>

            <div className="chat-mode-selector">
                <button
                    className={
                        chatMode === "knowledge"
                            ? "chat-mode-button active"
                            : "chat-mode-button"
                    }

                    disabled={isExistingSession}

                    onClick={() => setNewChatMode("knowledge")}

                >
                    Knowledge Base
                </button>

                <button
                    className={
                        chatMode === "note"
                            ? "chat-mode-button active"
                            : "chat-mode-button"
                    }

                    disabled={isExistingSession || !selectedNote}

                    onClick={() => setNewChatMode("note")}
                >
                    Current Note
                </button>
            </div>

            {chatMode === "note" ? (

                <div className="chat-context">
                    <span>Context:</span>

                    <strong>
                        {isExistingSession
                            ? currentSession?.title
                            : selectedNote?.title}
                    </strong>
                </div>

            ) : (

                <div className="chat-context">
                    <span>Context:</span>

                    <strong>Entire knowledge base</strong>
                </div>

            )}

            <Chat
                messages={messages}
                isSending={isSending}
                isLoadingMessages={isLoadingMessages}
                onSend={handleSend}
                mode={chatMode}
            />

        </section>
    );
}
