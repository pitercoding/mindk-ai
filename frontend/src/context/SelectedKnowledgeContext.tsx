import {
    createContext,
    useContext,
    useState,
    type ReactNode,
} from "react";

import type { Document } from "@/types/document";
import type { KnowledgeSelection } from "@/types/knowledge";
import type { Note } from "@/types/note";

interface SelectedKnowledgeContextType {
    selected: KnowledgeSelection | null;
    selectNote: (note: Note) => void;
    selectDocument: (document: Document) => void;
    clearSelection: () => void;
}

const SelectedKnowledgeContext = createContext<
    SelectedKnowledgeContextType | undefined
>(undefined);

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

export function useSelectedKnowledge() {

    const context = useContext(
        SelectedKnowledgeContext,
    );

    if (!context) {
        throw new Error(
            "useSelectedKnowledge must be used inside SelectedKnowledgeProvider"
        );
    }

    return context;
}
