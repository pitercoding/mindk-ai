import type { ChatHistory } from "@/types/chat-history";

interface ChatHistoryListProps {
    history: ChatHistory[];
    onSelect: (item: ChatHistory) => void;
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
    onSelect
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
                    onClick={() => onSelect(item)}
                    role="button"
                    tabIndex={0}
                    onKeyDown={(event) => {
                        if (event.key === "Enter" || event.key === " ") {
                            onSelect(item);
                        }
                    }}
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
