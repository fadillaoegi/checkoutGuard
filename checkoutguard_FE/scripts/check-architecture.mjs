import { readdirSync, readFileSync } from 'node:fs'
import { dirname, join, relative, resolve, sep } from 'node:path'
import ts from 'typescript'

const root = resolve('src')
const layers = ['domain', 'application', 'infrastructure', 'presentation']
const allowed = {
  domain: ['domain'],
  application: ['domain', 'application'],
  infrastructure: ['domain', 'application', 'infrastructure'],
  presentation: ['domain', 'application', 'presentation'],
}
const errors = []
function check(directory) {
  for (const entry of readdirSync(directory, { withFileTypes: true })) {
    const path = join(directory, entry.name)
    if (entry.isDirectory()) {
      check(path)
      continue
    }
    if (!/\.tsx?$/.test(entry.name)) continue
    const layer = relative(root, path).split(sep)[0]
    if (!layers.includes(layer)) continue // main.tsx is the composition root.
    const source = ts.createSourceFile(
      path,
      readFileSync(path, 'utf8'),
      ts.ScriptTarget.Latest,
      true,
    )
    const dependencies = []
    const visit = (node) => {
      if (
        (ts.isImportDeclaration(node) || ts.isExportDeclaration(node)) &&
        node.moduleSpecifier &&
        ts.isStringLiteral(node.moduleSpecifier)
      )
        dependencies.push(node.moduleSpecifier.text)
      if (
        ts.isCallExpression(node) &&
        (node.expression.kind === ts.SyntaxKind.ImportKeyword ||
          (ts.isIdentifier(node.expression) && node.expression.text === 'require'))
      ) {
        const value = node.arguments[0]
        if (value && ts.isStringLiteral(value)) dependencies.push(value.text)
        else errors.push(`${relative(root, path)} uses a nonliteral dynamic import`)
      }
      ts.forEachChild(node, visit)
    }
    visit(source)
    for (const dependency of dependencies) {
      if (!dependency.startsWith('.')) {
        if (layer === 'domain' || layer === 'application')
          errors.push(`${relative(root, path)} imports framework/library ${dependency}`)
        continue
      }
      const target = relative(root, resolve(dirname(path), dependency)).split(sep)[0]
      if (!allowed[layer].includes(target))
        errors.push(`${relative(root, path)} imports outward dependency ${dependency}`)
    }
  }
}
check(root)
if (errors.length) throw new Error(errors.join('\n'))
console.log(
  'Frontend dependency directions passed (domain → application ← adapters; composition in main.tsx)',
)
