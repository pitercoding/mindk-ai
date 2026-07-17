import type { Note } from "@/types/note";

interface KnowledgeBaseProps {
    notes: Note[];
    selectedNote: Note | null;
    onSelect: (note: Note) => void;
}

export default function KnowledgeBase({
    notes,
    selectedNote,
    onSelect,
}: KnowledgeBaseProps) {

    return (
        <section className="dashboard-panel">

            <header className="panel-header">

                <div>
                    <h2>KNOWLEDGE BASE</h2>
                    <span>Recent</span>
                </div>

                <button>
                    ...
                </button>

            </header>

            <div className="knowledge-list">

                {notes.map((note) => (

                    <article
                        key={note.id}
                        className={`knowledge-item ${selectedNote?.id === note.id
                                ? "knowledge-item-active"
                                : ""
                            }`}
                        onClick={() => onSelect(note)}
                    >

                        <h3>{note.title}</h3>

                        <p>
                            {note.content}
                        </p>

                    </article>

                ))}

            </div>

        </section>
    );
}
