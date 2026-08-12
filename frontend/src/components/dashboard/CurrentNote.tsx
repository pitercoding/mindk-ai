import MarkdownViewer from "@/components/common/MarkdownViewer";

import type { KnowledgeSelection } from "@/types/knowledge";

interface CurrentNoteProps {
    selection: KnowledgeSelection | null;
}

export default function CurrentNote({
    selection,
}: CurrentNoteProps) {

    const heading =
        selection?.type === "document"
            ? "CURRENT DOCUMENT"
            : "CURRENT NOTE";

    const name =
        selection?.type === "note"
            ? selection.item.title
            : selection?.item.name;

    const content =
        selection?.type === "note"
            ? stripTitleHeading(
                selection.item.content,
                selection.item.title,
            )
            : selection?.item.content ?? "";

    return (
        <section className="dashboard-panel current-note-panel">

            <header className="panel-header">

                <h2>{heading}</h2>

                <button>
                    ...
                </button>

            </header>

            <div className="current-note-content">

                {selection ? (
                    <>
                        <div className="current-note-title">
                            <h3>{name}</h3>

                            {selection.type === "document" && (
                                <span className="current-note-type">
                                    {selection.item.type}
                                </span>
                            )}
                        </div>

                        <div className="current-note-scroll">
                            <MarkdownViewer
                                content={content}
                            />
                        </div>
                    </>
                ) : (
                    <p className="current-note-empty">
                        Select a note or document to view content.
                    </p>
                )}

            </div>

        </section>
    );
}

function stripTitleHeading(content: string, title: string) {
    return content.replace(
        new RegExp(`^#\\s+${escapeRegExp(title)}\\s*\\n+`),
        "",
    );
}

function escapeRegExp(value: string) {
    return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
