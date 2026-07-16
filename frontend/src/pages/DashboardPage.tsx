import ChatPanel from "@/components/dashboard/ChatPanel";
import CurrentNote from "@/components/dashboard/CurrentNote";
import KnowledgeBase from "@/components/dashboard/KnowledgeBase";

export default function DashboardPage() {
    return (
        <section className="dashboard-page">

            <div className="dashboard-grid">

                <KnowledgeBase />

                <CurrentNote />

                <ChatPanel />

            </div>

        </section>
    );
}
