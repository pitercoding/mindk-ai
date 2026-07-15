import type { ChatHistory } from "@/types/chat-history";

interface ChatHistoryListProps {
    history: ChatHistory[];
}

const dateFormatter = new Intl.DateTimeFormat(undefined, {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
});

function formatHistoryDate(value: string) {
    const date = new Date(value);

    if (Number.isNaN(date.getTime())) {
        return "";
    }

    return dateFormatter.format(date);
}

export default function ChatHistoryList({
    history,
}: ChatHistoryListProps) {

    return (
        <section className="chat-history">

            <header className="history-header">
                <h2>Conversation History</h2>
            </header>

            {history.map((item) => (
                <article
                    key={item.id}
                    className="history-item"
                >
                    <p className="history-question">
                        {item.question}
                    </p>

                    <span className="history-date">
                        {formatHistoryDate(item.created_at)}
                    </span>

                </article>
            ))}

        </section>
    );
}
