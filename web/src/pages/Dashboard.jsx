import { useEffect, useState } from "react";
import { GraduationCap, Users, BookOpen, TrendingUp } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getStats } from "@/lib/api";

const statCards = [
  { key: "students", label: "Total Students", icon: GraduationCap, color: "text-blue-500" },
  { key: "teachers", label: "Total Teachers", icon: Users,          color: "text-violet-500" },
  { key: "classes",  label: "Total Classes",  icon: BookOpen,       color: "text-emerald-500" },
];

export default function Dashboard() {
  const [stats, setStats] = useState({ students: 0, teachers: 0, classes: 0 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getStats()
      .then(setStats)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold tracking-tight">Overview</h2>
        <p className="text-sm text-muted-foreground mt-1">
          Your school at a glance
        </p>
      </div>

      {/* Stat cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {statCards.map(({ key, label, icon: Icon, color }) => (
          <Card key={key}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {label}
              </CardTitle>
              <Icon className={`size-4 ${color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">
                {loading ? (
                  <span className="inline-block h-8 w-12 animate-pulse rounded bg-muted" />
                ) : (
                  stats[key] ?? 0
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Info banner */}
      <Card>
        <CardContent className="flex items-center gap-4 py-5">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-primary/10">
            <TrendingUp className="size-5 text-primary" />
          </div>
          <div>
            <p className="text-sm font-medium">Getting started</p>
            <p className="text-sm text-muted-foreground">
              Add teachers first, then create classes and assign them. Finally enrol students into classes.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
