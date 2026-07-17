import type { Note } from "@/types/note";

interface CurrentNoteProps {
    note: Note | null;
}

export default function CurrentNote({
    note,
}: CurrentNoteProps) {

    return (
        <section className="dashboard-panel">

            <header className="panel-header">

                <h2>CURRENT NOTE</h2>

                <button>
                    ...
                </button>

            </header>

            <div>

                {note ? (
                    <>
                        <h3>{note.title}</h3>

                        <p>{note.content}</p>
                    </>
                ) : (
                    <p>
                        Select a note to view content.
                    </p>
                )}

            </div>

        </section>
    );
}
