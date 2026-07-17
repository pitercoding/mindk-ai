import MarkdownViewer from "@/components/common/MarkdownViewer";

import type { Note } from "@/types/note";

interface CurrentNoteProps {
    note: Note | null;
}

export default function CurrentNote({
    note,
}: CurrentNoteProps) {

    const content = note?.content.replace(
        new RegExp(`^#\\s+${escapeRegExp(note.title)}\\s*\\n+`),
        "",
    ) ?? "";

    return (
        <section className="dashboard-panel current-note-panel">

            <header className="panel-header">

                <h2>CURRENT NOTE</h2>

                <button>
                    ...
                </button>

            </header>

            <div className="current-note-content">

                {note ? (
                    <>
                        <div className="current-note-title">
                            <h3>{note.title}</h3>
                        </div>

                        <div className="current-note-scroll">
                            <MarkdownViewer
                                content={content ?? ""}
                            />
                        </div>
                    </>
                ) : (
                    <p className="current-note-empty">
                        Select a note to view content.
                    </p>
                )}

            </div>

        </section>
    );
}

function escapeRegExp(value: string) {
    return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
