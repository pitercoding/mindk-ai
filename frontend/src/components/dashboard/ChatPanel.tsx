import { useSelectedNote } from "@/context/SelectedNoteContext";

export default function ChatPanel() {

    const {
        selectedNote,
    } = useSelectedNote();

    return (
        <section className="dashboard-panel chat-panel">

            <header className="panel-header">

                <div>
                    <h2>MINDK CHAT</h2>

                    <span>
                        Ask MindK anything about your data
                    </span>
                </div>

                <button>
                    ...
                </button>

            </header>


            <div className="chat-empty-state">

                {selectedNote && (
                    <div className="chat-context">

                        <span>
                            Context:
                        </span>

                        <strong>
                            {selectedNote.title}
                        </strong>

                    </div>
                )}

                <p>
                    Chat messages will appear here.
                </p>

            </div>

        </section>
    );
}
