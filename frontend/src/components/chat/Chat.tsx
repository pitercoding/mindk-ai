import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";
import EmptyState from "@/components/common/EmptyState";
import { SparkIcon } from "@/components/common/icons";
import { RowSkeletons } from "@/components/common/Skeleton";

import type { Message } from "@/types/message";

interface ChatProps {
    messages: Message[];
    isSending: boolean;
    isLoadingMessages: boolean;
    onSend: (message: string) => Promise<void>;
    mode: "knowledge" | "note";
}

const SUGGESTIONS: Record<ChatProps["mode"], string[]> = {
    knowledge: [
        "Summarize my latest notes",
        "Which documents have I uploaded?",
        "List the action items across my notes",
    ],
    note: [
        "Summarize this note",
        "List the key points",
        "What questions does this note raise?",
    ],
};

export default function Chat({
    messages,
    isSending,
    isLoadingMessages,
    onSend,
    mode,
}: ChatProps) {

    const emptyTitle =
        mode === "knowledge"
            ? "Ask anything about your knowledge base"
            : "Ask MindK about this note";

    const emptyDescription =
        mode === "knowledge"
            ? "Search across your notes and documents using AI."
            : "Ask questions about the selected note.";

    const isBusy = isSending || isLoadingMessages;

    return (
        <>
            <section className="chat-content">

                {isLoadingMessages ? (

                    <RowSkeletons
                        count={3}
                        label="Loading messages"
                        className="chat-skeleton"
                    />

                ) : messages.length === 0 ? (

                    <EmptyState
                        icon={<SparkIcon size={20} />}
                        title={emptyTitle}
                        description={emptyDescription}
                        actions={
                            <ul className="suggestion-list" aria-label="Suggested questions">
                                {SUGGESTIONS[mode].map((suggestion) => (
                                    <li key={suggestion}>
                                        <button
                                            type="button"
                                            className="suggestion-chip"
                                            disabled={isBusy}
                                            onClick={() => onSend(suggestion)}
                                        >
                                            {suggestion}
                                        </button>
                                    </li>
                                ))}
                            </ul>
                        }
                    />

                ) : (

                    <ChatMessageList
                        messages={messages}
                        isSending={isSending}
                    />

                )}

            </section>

            <footer className="chat-footer">

                <ChatInput
                    onSend={onSend}
                    disabled={isBusy}
                />

                <p className="chat-hint">
                    Enter to send · Shift+Enter for a new line
                </p>

            </footer>
        </>
    );
}
