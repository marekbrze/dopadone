#!/bin/bash
# Seed database with test data for manual testing

set -e

DB_PATH="${1:-./dopa.db}"
CMD="go run ./cmd/dopa --db $DB_PATH"

echo "Seeding database at $DB_PATH..."
echo ""

# Reset database
echo "Resetting database..."
rm -f "$DB_PATH"
$CMD migrate up > /dev/null 2>&1
echo "  ✓ Database reset and migrations applied"
echo ""

# Create 3 areas
echo "Creating areas..."
AREA1=$($CMD areas create --name "Personal" --color "#3498db" --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
AREA2=$($CMD areas create --name "Work" --color "#e74c3c" --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
AREA3=$($CMD areas create --name "Side Projects" --color "#2ecc71" --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
echo "  ✓ Created 3 areas"

# Create subareas for each area
echo "Creating subareas..."
# Personal subareas
SUB1=$($CMD subareas create --name "Health & Fitness" --area-id $AREA1 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
SUB2=$($CMD subareas create --name "Learning" --area-id $AREA1 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
SUB3=$($CMD subareas create --name "Home" --area-id $AREA1 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')

# Work subareas
SUB4=$($CMD subareas create --name "Client A" --area-id $AREA2 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
SUB5=$($CMD subareas create --name "Client B" --area-id $AREA2 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')

# Side Projects subareas
SUB6=$($CMD subareas create --name "Open Source" --area-id $AREA3 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
SUB7=$($CMD subareas create --name "Mobile App" --area-id $AREA3 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
echo "  ✓ Created 7 subareas"

# Create projects with hierarchy
echo "Creating projects..."
# Personal - Health & Fitness
PROJ1=$($CMD projects create --name "Marathon Training" --subarea-id $SUB1 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
PROJ2=$($CMD projects create --name "Meal Prep System" --subarea-id $SUB1 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
$CMD projects create --name "Week 1-4 Base Building" --parent-id $PROJ1 --output json > /dev/null
$CMD projects create --name "Week 5-8 Speed Work" --parent-id $PROJ1 --output json > /dev/null
$CMD projects create --name "Week 9-12 Long Runs" --parent-id $PROJ1 --output json > /dev/null

# Personal - Learning
PROJ3=$($CMD projects create --name "Learn Rust" --subarea-id $SUB2 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
$CMD projects create --name "Rust Book Chapters 1-10" --parent-id $PROJ3 --output json > /dev/null
$CMD projects create --name "Rust Book Chapters 11-20" --parent-id $PROJ3 --output json > /dev/null

PROJ4=$($CMD projects create --name "AWS Certification" --subarea-id $SUB2 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
$CMD projects create --name "SAA-C02 Exam Prep" --parent-id $PROJ4 --output json > /dev/null

# Personal - Home
$CMD projects create --name "Garden Renovation" --subarea-id $SUB3 --output json > /dev/null
$CMD projects create --name "Basement Organization" --subarea-id $SUB3 --output json > /dev/null

# Work - Client A
PROJ5=$($CMD projects create --name "E-commerce Platform" --subarea-id $SUB4 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
$CMD projects create --name "Backend API" --parent-id $PROJ5 --output json > /dev/null
$CMD projects create --name "Frontend UI" --parent-id $PROJ5 --output json > /dev/null
$CMD projects create --name "Payment Integration" --parent-id $PROJ5 --output json > /dev/null

$CMD projects create --name "Mobile App v2.0" --subarea-id $SUB4 --output json > /dev/null

# Work - Client B
$CMD projects create --name "Analytics Dashboard" --subarea-id $SUB5 --output json > /dev/null
$CMD projects create --name "Data Pipeline" --subarea-id $SUB5 --output json > /dev/null

# Side Projects - Open Source
PROJ6=$($CMD projects create --name "Dopadone" --subarea-id $SUB6 --output json | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"$' | tr -d '"')
$CMD projects create --name "Core Features" --parent-id $PROJ6 --output json > /dev/null
$CMD projects create --name "TUI Interface" --parent-id $PROJ6 --output json > /dev/null
$CMD projects create --name "CLI Commands" --parent-id $PROJ6 --output json > /dev/null

# Side Projects - Mobile App
$CMD projects create --name "Habit Tracker" --subarea-id $SUB7 --output json > /dev/null
$CMD projects create --name "Note Taking App" --subarea-id $SUB7 --output json > /dev/null
echo "  ✓ Created 25 projects (with hierarchy)"

# Create tasks for various projects
echo "Creating tasks..."
# Get project IDs and names
PROJ_DATA=$($CMD projects list | awk 'NR>1 {print $1 "|" $2}')

# Create different tasks based on project type
COUNT=0
echo "$PROJ_DATA" | while IFS='|' read PROJ_ID PROJ_NAME; do
    if [ -n "$PROJ_ID" ]; then
        case "$PROJ_NAME" in
            *Training*|*Running*|*Base*|*Speed*|*Long*)
                # Fitness/Training tasks
                $CMD tasks create --title "Plan workout schedule" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Buy running gear" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Track progress" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *Meal*|*Prep*)
                # Meal prep tasks
                $CMD tasks create --title "Plan weekly menu" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Create shopping list" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Prep ingredients" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *Rust*|*AWS*|*Learn*|*Book*|*Exam*)
                # Learning tasks
                $CMD tasks create --title "Read documentation" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Complete exercises" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Take practice tests" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *Garden*|*Renovation*)
                # Home improvement tasks
                $CMD tasks create --title "Research materials" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Get quotes from contractors" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Set budget" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *Basement*|*Organization*)
                # Organization tasks
                $CMD tasks create --title "Declutter items" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Buy storage containers" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Organize by category" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *API*|*Backend*|*Frontend*|*UI*|*CLI*|*Payment*|*Core*|*Features*|*Interface*)
                # Software development tasks
                $CMD tasks create --title "Define requirements" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Write code" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Write tests" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *E-commerce*|*Analytics*|*Dashboard*|*Platform*|*Pipeline*)
                # Business/Analytics tasks
                $CMD tasks create --title "Gather requirements" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Design system architecture" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Implement MVP" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *Habit*|*Tracker*|*Note*|*App*|*Mobile*)
                # Mobile app tasks
                $CMD tasks create --title "Design UI mockups" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Implement core functionality" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Test on devices" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *Dopadone*|*Open*|*Source*)
                # Open source tasks
                $CMD tasks create --title "Write documentation" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Set up CI/CD" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Create release" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
            *)
                # Generic personal tasks
                $CMD tasks create --title "Research options" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Create action plan" --project-id "$PROJ_ID" > /dev/null 2>&1
                $CMD tasks create --title "Execute plan" --project-id "$PROJ_ID" > /dev/null 2>&1
                ;;
        esac
        COUNT=$((COUNT + 3))
    fi
done
echo "  ✓ Created contextual tasks"
echo ""
echo "✅ Database seeded successfully!"
echo ""
echo "Summary:"
$CMD areas list
echo ""
echo "Run 'go run ./cmd/dopa tui' to test the TUI with this data."
