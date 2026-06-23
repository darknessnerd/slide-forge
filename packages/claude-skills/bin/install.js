#!/usr/bin/env node
// claude-skills install
// Copies skills from:
//   1. packages/claude-skills/skills/     (team-authored skills, flat .md files)
//   2. node_modules/caveman-installer/skills/<name>/SKILL.md  (JuliusBrussee/caveman)
// into .claude/skills/ in the consumer repo root.
//
// Usage:
//   npm run setup:claude           normal install
//   npm run setup:claude -- --dry-run   preview without writing

const fs = require('fs')
const path = require('path')

const dryRun = process.argv.includes('--dry-run')

const targetDir = path.join(process.cwd(), '.claude', 'skills')

if (!dryRun) {
  fs.mkdirSync(targetDir, { recursive: true })
}

// ── Helpers ────────────────────────────────────────────────────────────────

function resolveNodeModules(pkgName) {
  const candidates = [
    path.join(__dirname, '..', 'node_modules', pkgName),
    path.join(__dirname, '..', '..', 'node_modules', pkgName),
    path.join(process.cwd(), 'node_modules', pkgName),
  ]
  return candidates.find(fs.existsSync) || null
}

// Install flat .md files from a directory
function installFlatSkills(label, dir) {
  if (!dir || !fs.existsSync(dir)) {
    console.warn(`WARN: skills dir not found for ${label} — skipping`)
    return 0
  }
  const files = fs.readdirSync(dir).filter(f => f.endsWith('.md'))
  for (const file of files) {
    copySkill(label, path.join(dir, file), file)
  }
  return files.length
}

// Install skills from caveman-installer layout: skills/<name>/SKILL.md → <name>.md
function installCavemanSkills(label, pkgDir) {
  if (!pkgDir) {
    console.warn(`WARN: package not found for ${label} — skipping`)
    return 0
  }
  const skillsDir = path.join(pkgDir, 'skills')
  if (!fs.existsSync(skillsDir)) {
    console.warn(`WARN: skills dir not found in ${label} — skipping`)
    return 0
  }
  let count = 0
  for (const name of fs.readdirSync(skillsDir)) {
    const skillFile = path.join(skillsDir, name, 'SKILL.md')
    if (fs.existsSync(skillFile)) {
      copySkill(label, skillFile, `${name}.md`)
      count++
    }
  }
  return count
}

function copySkill(label, src, destName) {
  const dest = path.join(targetDir, destName)
  if (dryRun) {
    console.log(`[dry-run] ${label}/${destName} → .claude/skills/${destName}`)
  } else {
    fs.copyFileSync(src, dest)
    console.log(`installed [${label}]: .claude/skills/${destName}`)
  }
}

// ── Install ────────────────────────────────────────────────────────────────

let total = 0

// 1. Team skills (flat .md files)
total += installFlatSkills(
  '@team/claude-skills',
  path.join(__dirname, '..', 'skills')
)

// 2. caveman-installer (JuliusBrussee/caveman)
total += installCavemanSkills(
  'caveman-installer',
  resolveNodeModules('caveman-installer')
)

console.log(dryRun
  ? `\nDry run — ${total} file(s) would be installed.`
  : `\nDone — ${total} skill(s) installed.`)
