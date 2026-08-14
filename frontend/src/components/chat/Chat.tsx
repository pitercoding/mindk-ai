import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";

import type { Message } from "@/types/message";

interface ChatProps {
    messages: Message[];
    isSending: boolean;
    isLoadingMessages: boolean;
    onSend: (message: string) => Promise<void>;
    mode: "knowledge" | "note";
}

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

    return (
        <>
            <section className="chat-content">

                {isLoadingMessages ? (

                    <div className="chat-empty-state">

                        <h3>Loading messages...</h3>

                    </div>

                ) : messages.length === 0 ? (

                    <div className="chat-empty-state">

                        <h3>{emptyTitle}</h3>

                        <p>{emptyDescription}</p>

                    </div>

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
                    disabled={isSending || isLoadingMessages}
                />

            </footer>
        </>
    );
}
