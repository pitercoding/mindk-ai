import {
    useState,
    type ReactNode,
} from "react";

import { SelectedKnowledgeContext } from "@/context/selected-knowledge-context";

import type { Document } from "@/types/document";
import type { KnowledgeSelection } from "@/types/knowledge";
import type { Note } from "@/types/note";

interface SelectedKnowledgeProviderProps {
    children: ReactNode;
}

export function SelectedKnowledgeProvider({
    children,
}: SelectedKnowledgeProviderProps) {

    const [selected, setSelected] =
        useState<KnowledgeSelection | null>(null);

    function selectNote(note: Note) {
        setSelected({ type: "note", item: note });
    }

    function selectDocument(document: Document) {
        setSelected({ type: "document", item: document });
    }

    function clearSelection() {
        setSelected(null);
    }

    return (
        <SelectedKnowledgeContext.Provider
            value={{
                selected,
                selectNote,
                selectDocument,
                clearSelection,
            }}
        >
            {children}
        </SelectedKnowledgeContext.Provider>
    );
}
