import { useEffect, useState } from "react";

import Chat from "@/components/chat/Chat";

import { useSelectedNote } from "@/context/SelectedNoteContext";

import { sendMessage } from "@/services/chatService";
import { getMessagesByNote } from "@/services/chatMessageService";

import type { Message } from "@/types/message";


type ChatMode = "note" | "knowledge";


export default function ChatPanel() {

    const [chatMode, setChatMode] = useState<ChatMode>("knowledge");

    const [chatSessions, setChatSessions] = useState<Record<number, Message[]>>({});

    const [knowledgeMessages, setKnowledgeMessages] = useState<Message[]>([]);

    const [isLoading, setIsLoading] = useState(false);

    const { selectedNote } = useSelectedNote();

    const noteMessages =
        selectedNote
            ? chatSessions[selectedNote.id] ?? []
            : [];

    const messages =
        chatMode === "note"
            ? noteMessages
            : knowledgeMessages;


    useEffect(() => {
        if (chatMode !== "note" || !selectedNote) {
            return;
        }

        const note = selectedNote;

        async function loadMessages() {

            try {

                const history = (await getMessagesByNote(note.id)) ?? [];

                setChatSessions(
                    (previousSessions) => ({
                        ...previousSessions,
                        [note.id]:

                            history.length > 0
                                ? history.map(message => ({
                                    id: message.id,
                                    role: message.role,
                                    content: message.content,
                                }))
                                : []
                    }),
                );

            } catch (error) {

                console.error("Failed to load chat messages:", error);
            }
        }

        loadMessages();

    }, [selectedNote, chatMode]);


    function updateCurrentMessages(
        updater: (messages: Message[]) => Message[],
    ) {

        if (chatMode === "knowledge") {

            setKnowledgeMessages(previous =>
                updater(previous),
            );

            return;
        }

        if (!selectedNote) {
            return;
        }

        setChatSessions(previousSessions => {

            const currentMessages =
                previousSessions[selectedNote.id] ?? [];

            return {
                ...previousSessions,

                [selectedNote.id]:
                    updater(currentMessages),
            };
        });
    }

    async function handleSend(
        message: string,
    ) {

        if (
            chatMode === "note" &&
            selectedNote &&
            messages.length === 0
        ) {

            updateCurrentMessages(previousMessages => [

                ...previousMessages,

                {
                    id: crypto.randomUUID(),
                    role: "assistant",
                    content:
                        `I am ready to answer questions about "${selectedNote.title}".\n\n` +
                        `You can ask me to:\n\n` +
                        `- summarize this note\n` +
                        `- explain concepts\n` +
                        `- provide examples\n` +
                        `- clarify information`,
                },
            ]);
        }

        const userMessage: Message = {
            id: crypto.randomUUID(),
            role: "user",
            content: message,
        };

        updateCurrentMessages(
            previousMessages => [
                ...previousMessages,
                userMessage,
            ],
        );

        setIsLoading(true);

        try {

            const loadingId = crypto.randomUUID();

            updateCurrentMessages(
                previousMessages => [
                    ...previousMessages,

                    {
                        id: loadingId,
                        role: "assistant",
                        content:
                            "MindK AI is analyzing your knowledge...",
                    },
                ],
            );

            const response =
                await sendMessage({

                    message,

                    context:
                        chatMode === "note" && selectedNote
                            ? {
                                note_id: selectedNote.id,
                                title: selectedNote.title,
                                content: selectedNote.content,
                            }
                            : undefined,

                });

            updateCurrentMessages(
                previousMessages =>
                    previousMessages.filter(
                        item =>
                            item.id !== loadingId,
                    ),
            );

            updateCurrentMessages(
                previousMessages => [

                    ...previousMessages,

                    {
                        id: crypto.randomUUID(),
                        role: "assistant",
                        content: response.answer,
                    },

                ],
            );

        } catch (error) {

            console.error(
                "Failed to send dashboard chat message:",
                error,
            );

            updateCurrentMessages(
                previousMessages => [

                    ...previousMessages,

                    {
                        id: crypto.randomUUID(),
                        role: "assistant",
                        content:
                            "Sorry, something went wrong.",
                    },

                ],
            );

        } finally {

            setIsLoading(false);
        }
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

                    onClick={() => setChatMode("knowledge")}

                >
                    Knowledge Base
                </button>

                <button
                    className={
                        chatMode === "note"
                            ? "chat-mode-button active"
                            : "chat-mode-button"
                    }

                    disabled={!selectedNote}

                    onClick={() => setChatMode("note")}
                >
                    Current Note
                </button>
            </div>

            {chatMode === "note" && selectedNote && (

                <div className="chat-context">
                    <span>Context:</span>

                    <strong>{selectedNote.title}</strong>
                </div>

            )}

            {chatMode === "knowledge" && (

                <div className="chat-context">
                    <span>Context:</span>

                    <strong>Entire knowledge base</strong>
                </div>

            )}

            <Chat
                messages={messages}
                isLoading={isLoading}
                onSend={handleSend}
                mode={chatMode}
            />

        </section>
    );
}
