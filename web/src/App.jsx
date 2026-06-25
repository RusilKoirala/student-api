import { useEffect, useState } from "react";
import { Toaster } from "@/components/ui/sonner";
import Layout from "@/components/Layout";
import AuthPage from "@/pages/AuthPage";
import Dashboard from "@/pages/Dashboard";
import Students from "@/pages/Students";
import Teachers from "@/pages/Teachers";
import Classes from "@/pages/Classes";

const pages = {
  dashboard: Dashboard,
  students:  Students,
  teachers:  Teachers,
  classes:   Classes,
};

function getInitialAuth() {
  const token  = localStorage.getItem("token");
  const school = localStorage.getItem("school");
  if (token && school) {
    try { return { token, ...JSON.parse(school) }; }
    catch { /* corrupted — fall through */ }
  }
  return null;
}

export default function App() {
  const [auth, setAuth] = useState(getInitialAuth);
  const [page, setPage] = useState("dashboard");

  // Listen for 401 auto-logout fired from api.js interceptor
  useEffect(() => {
    const handle = () => setAuth(null);
    window.addEventListener("auth:logout", handle);
    return () => window.removeEventListener("auth:logout", handle);
  }, []);

  function handleAuth(data) {
    setAuth({ token: data.token, name: data.schoolName, id: data.schoolId });
  }

  function handleLogout() {
    localStorage.removeItem("token");
    localStorage.removeItem("school");
    setAuth(null);
    setPage("dashboard");
  }

  if (!auth) {
    return (
      <>
        <AuthPage onAuth={handleAuth} />
        <Toaster richColors position="bottom-right" />
      </>
    );
  }

  const Page = pages[page] ?? Dashboard;

  return (
    <>
      <Layout page={page} setPage={setPage} school={auth} onLogout={handleLogout}>
        <Page />
      </Layout>
      <Toaster richColors position="bottom-right" />
    </>
  );
}
