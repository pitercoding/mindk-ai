import { Link } from "react-router-dom";

import EmptyState from "@/components/common/EmptyState";
import { SearchIcon } from "@/components/common/icons";

export default function NotFoundPage() {
    return (
        <section className="not-found-page">

            <EmptyState
                icon={<SearchIcon size={20} />}
                title="Page not found"
                description="The page you are looking for does not exist or was moved."
                actions={
                    <Link to="/" className="btn btn-primary">
                        Back to dashboard
                    </Link>
                }
            />

        </section>
    );
}
