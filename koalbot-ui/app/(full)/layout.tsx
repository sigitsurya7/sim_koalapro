export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
    return(
        <section className="p-4">
            {children}
        </section>
    )
}