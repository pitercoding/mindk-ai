import { useEffect, useState } from "react";

import ChatPanel from "@/components/dashboard/ChatPanel";
import CurrentNote from "@/components/dashboard/CurrentNote";
import KnowledgeBase from "@/components/dashboard/KnowledgeBase";

import { useSelectedKnowledge } from "@/context/SelectedKnowledgeContext";

import { getDocuments } from "@/services/documentService";
import { getNotes } from "@/services/noteService";

import type { Document } from "@/types/document";
import type { Note } from "@/types/note";

type DashboardTab = "knowledge" | "note" | "chat";

const dashboardTabs: { id: DashboardTab; label: string }[] = [
    { id: "knowledge", label: "Knowledge" },
    { id: "note", label: "Note" },
    { id: "chat", label: "Chat" },
];

export default function DashboardPage() {

    const [notes, setNotes] = useState<Note[]>([]);
    const [documents, setDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [activeTab, setActiveTab] = useState<DashboardTab>("knowledge");

    const {
        selected,
        selectNote,
        selectDocument,
        clearSelection,
    } = useSelectedKnowledge();


    async function loadKnowledgeBase() {

        setIsLoading(true);
        setError(null);

        try {

            const [notesResponse, documentsResponse] = await Promise.all([
                getNotes(),
                getDocuments(),
            ]);

            setNotes(notesResponse);
            setDocuments(documentsResponse);

            if (selected?.type === "note") {

                const freshNote = notesResponse.find(
                    (note) => note.id === selected.item.id
                );

                if (freshNote) {
                    selectNote(freshNote);
                } else {
                    clearSelection();
                }

            } else if (selected?.type === "document") {

                const freshDocument = documentsResponse.find(
                    (document) => document.id === selected.item.id
                );

                if (freshDocument) {
                    selectDocument(freshDocument);
                } else {
                    clearSelection();
                }

            } else if (notesResponse.length > 0) {
                selectNote(notesResponse[0]);
            } else if (documentsResponse.length > 0) {
                selectDocument(documentsResponse[0]);
            }

        } catch (error) {

            console.error(
                "Failed to load knowledge base:",
                error,
            );

            setError("Failed to load knowledge base.");

        } finally {
            setIsLoading(false);
        }
    }


    useEffect(() => {
        loadKnowledgeBase();
    }, []);


    return (
        <section className="dashboard-page">

            {/* Mobile only: one panel at a time (hidden by CSS on larger screens) */}
            <div className="dashboard-tabs" role="tablist" aria-label="Dashboard sections">
                {dashboardTabs.map((tab) => (
                    <button
                        key={tab.id}
                        type="button"
                        role="tab"
                        aria-selected={activeTab === tab.id}
                        className={
                            activeTab === tab.id
                                ? "dashboard-tab active"
                                : "dashboard-tab"
                        }
                        onClick={() => setActiveTab(tab.id)}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>

            <div className="dashboard-grid" data-active-tab={activeTab}>

                <KnowledgeBase
                    notes={notes}
                    documents={documents}
                    isLoading={isLoading}
                    error={error}
                    onItemSelected={() => setActiveTab("note")}
                />

                <CurrentNote
                    selection={selected}
                />

                <ChatPanel />

            </div>

        </section>
    );
}
