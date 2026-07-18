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
        <section className="dashboard-panel knowledge-panel">

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

                {notes.map((note) => {

                    function handleSelect() {
                        onSelect(note);
                    }

                    return (
                        <article
                            key={note.id}
                            className={`knowledge-item ${selectedNote?.id === note.id
                                    ? "knowledge-item-active"
                                    : ""
                                }`}
                            role="button"
                            tabIndex={0}
                            onClick={handleSelect}
                            onKeyDown={(event) => {
                                if (
                                    event.key === "Enter" ||
                                    event.key === " "
                                ) {
                                    event.preventDefault();
                                    handleSelect();
                                }
                            }}
                        >

                            <h3>{note.title}</h3>

                            <p>
                                {note.content}
                            </p>

                        </article>
                    );
                })}

            </div>

        </section>
    );
}
