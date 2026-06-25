import { useState } from "react";
import { School } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { register, login } from "@/lib/api";

export default function AuthPage({ onAuth }) {
  const [mode, setMode] = useState("login"); // "login" | "register"
  const [loading, setLoading] = useState(false);

  const [form, setForm] = useState({
    schoolName: "",
    username: "",
    password: "",
  });

  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.value }));

  async function handleSubmit(e) {
    e.preventDefault();
    setLoading(true);
    try {
      let data;
      if (mode === "register") {
        if (!form.schoolName) { toast.error("School name is required"); return; }
        data = await register({
          schoolName: form.schoolName,
          username:   form.username,
          password:   form.password,
        });
        toast.success(`Welcome, ${data.schoolName}!`);
      } else {
        data = await login({ username: form.username, password: form.password });
        toast.success(`Welcome back, ${data.schoolName}!`);
      }

      localStorage.setItem("token", data.token);
      localStorage.setItem("school", JSON.stringify({ name: data.schoolName, id: data.schoolId }));
      onAuth(data);
    } catch (err) {
      const msg = err.response?.data?.error || "Something went wrong";
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-svh flex items-center justify-center bg-background p-4">
      <div className="w-full max-w-sm space-y-6">
        {/* Logo */}
        <div className="flex flex-col items-center gap-3">
          <div className="flex size-12 items-center justify-center rounded-xl bg-primary">
            <School className="size-6 text-primary-foreground" />
          </div>
          <div className="text-center">
            <h1 className="text-2xl font-semibold tracking-tight">SchoolOS</h1>
            <p className="text-sm text-muted-foreground mt-1">
              {mode === "login" ? "Sign in to your school" : "Create a new school"}
            </p>
          </div>
        </div>

        <Card>
          <CardHeader className="pb-4">
            <CardTitle className="text-base">
              {mode === "login" ? "Sign In" : "Register School"}
            </CardTitle>
            <CardDescription>
              {mode === "login"
                ? "Enter your credentials to access your dashboard"
                : "Set up your school in seconds"}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {mode === "register" && (
                <div className="space-y-1.5">
                  <Label htmlFor="schoolName">School Name</Label>
                  <Input
                    id="schoolName"
                    placeholder="Sunrise Academy"
                    value={form.schoolName}
                    onChange={set("schoolName")}
                    required
                  />
                </div>
              )}

              <div className="space-y-1.5">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  placeholder="admin"
                  value={form.username}
                  onChange={set("username")}
                  autoComplete="username"
                  required
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  value={form.password}
                  onChange={set("password")}
                  autoComplete={mode === "login" ? "current-password" : "new-password"}
                  required
                />
              </div>

              <Button type="submit" className="w-full" disabled={loading}>
                {loading
                  ? mode === "login" ? "Signing in…" : "Creating school…"
                  : mode === "login" ? "Sign In" : "Create School"}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Toggle */}
        <p className="text-center text-sm text-muted-foreground">
          {mode === "login" ? (
            <>
              Don&apos;t have a school yet?{" "}
              <button
                onClick={() => setMode("register")}
                className="font-medium text-foreground underline-offset-4 hover:underline"
              >
                Create one
              </button>
            </>
          ) : (
            <>
              Already have an account?{" "}
              <button
                onClick={() => setMode("login")}
                className="font-medium text-foreground underline-offset-4 hover:underline"
              >
                Sign in
              </button>
            </>
          )}
        </p>
      </div>
    </div>
  );
}
