import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";
import type { ChatContext } from "@/types/chat";

import type { Message } from "@/types/message";

interface ChatProps {
    messages: Message[];
    isLoading: boolean;
    onSend: (message: string) => Promise<void>;
    context?: ChatContext;
}

export default function Chat({
    messages,
    isLoading,
    onSend,
}: ChatProps) {

    return (

        <>
            <section className="chat-content">

                {messages.length === 0 ? (

                    <div className="chat-empty-state">

                        <h3>Ask MindK about your knowledge</h3>

                        <p>Select a note and start asking questions.</p>
                    </div>

                ) : (

                    <ChatMessageList
                        messages={messages}
                    />
                )}

            </section>

            <footer className="chat-footer">

                <ChatInput
                    onSend={onSend}
                    disabled={isLoading}
                />

            </footer>
        </>
    );
}
