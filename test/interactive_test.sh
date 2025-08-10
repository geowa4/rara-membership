#!/bin/bash

# Interactive API Test Client
# This script provides an interactive menu for testing API endpoints

source "$(dirname "$0")/test_helpers.sh"

# Function to show menu
show_menu() {
    echo ""
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}  RARA API Interactive Tester${NC}"
    echo -e "${BLUE}================================${NC}"
    echo ""
    echo "1. List all members"
    echo "2. Create a member"
    echo "3. Update a member"
    echo "4. Get member points"
    echo "5. List all events"
    echo "6. Create an event"
    echo "7. Update an event"
    echo "8. Allocate points"
    echo "9. Redeem points"
    echo "10. Check balance"
    echo "11. Run automated tests"
    echo "0. Exit"
    echo ""
}

# Function to list members
list_members() {
    print_section "List All Members"
    response=$(curl -s "$API_URL/members")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to create member interactively
create_member_interactive() {
    print_section "Create New Member"
    
    read -p "Name: " name
    read -p "Call Sign: " call_sign
    read -p "Email: " email
    read -p "Phone (optional): " phone
    read -p "License Class (Technician/General/Extra): " license_class
    
    local data="{
        \"name\": \"$name\",
        \"call_sign\": \"$call_sign\",
        \"email\": \"$email\""
    
    [ ! -z "$phone" ] && data="$data, \"phone\": \"$phone\""
    [ ! -z "$license_class" ] && data="$data, \"license_class\": \"$license_class\""
    
    data="$data }"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/members/create")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to update member
update_member_interactive() {
    print_section "Update Member"
    
    read -p "Call Sign of member to update: " call_sign
    read -p "New Phone (press enter to skip): " phone
    read -p "New Email (press enter to skip): " email
    read -p "New License Class (press enter to skip): " license_class
    
    local data="{"
    local first=true
    
    if [ ! -z "$phone" ]; then
        data="$data \"phone\": \"$phone\""
        first=false
    fi
    
    if [ ! -z "$email" ]; then
        [ "$first" = false ] && data="$data,"
        data="$data \"email\": \"$email\""
        first=false
    fi
    
    if [ ! -z "$license_class" ]; then
        [ "$first" = false ] && data="$data,"
        data="$data \"license_class\": \"$license_class\""
    fi
    
    data="$data }"
    
    response=$(curl -s -X PUT \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/members/update?call_sign=$call_sign")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to get member points
get_member_points() {
    print_section "Get Member Points History"
    
    read -p "Call Sign: " call_sign
    
    response=$(curl -s "$API_URL/members/points?call_sign=$call_sign")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to list events
list_events() {
    print_section "List All Events"
    response=$(curl -s "$API_URL/events")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to create event
create_event_interactive() {
    print_section "Create New Event"
    
    read -p "Event Name: " name
    read -p "Description: " description
    read -p "Default Points (1-100): " points
    
    local data="{
        \"name\": \"$name\",
        \"description\": \"$description\",
        \"default_points\": $points
    }"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/events/create")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to update event
update_event_interactive() {
    print_section "Update Event"
    
    read -p "Event ID to update: " event_id
    read -p "New Name (press enter to skip): " name
    read -p "New Description (press enter to skip): " description
    read -p "New Default Points (press enter to skip): " points
    
    local data="{"
    local first=true
    
    if [ ! -z "$name" ]; then
        data="$data \"name\": \"$name\""
        first=false
    fi
    
    if [ ! -z "$description" ]; then
        [ "$first" = false ] && data="$data,"
        data="$data \"description\": \"$description\""
        first=false
    fi
    
    if [ ! -z "$points" ]; then
        [ "$first" = false ] && data="$data,"
        data="$data \"default_points\": $points"
    fi
    
    data="$data }"
    
    response=$(curl -s -X PUT \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/events/update?id=$event_id")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to allocate points
allocate_points_interactive() {
    print_section "Allocate Points"
    
    read -p "Event ID: " event_id
    read -p "Member ID: " member_id
    read -p "Points (press enter for event default): " points
    read -p "Notes (optional): " notes
    
    local data="{
        \"event_id\": $event_id,
        \"member_id\": $member_id"
    
    [ ! -z "$points" ] && data="$data, \"points\": $points"
    [ ! -z "$notes" ] && data="$data, \"notes\": \"$notes\""
    
    data="$data }"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/points/allocate")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to redeem points
redeem_points_interactive() {
    print_section "Redeem Points"
    
    read -p "Member ID: " member_id
    read -p "Points to redeem: " points
    read -p "Notes (what was redeemed): " notes
    
    local data="{
        \"member_id\": $member_id,
        \"points\": $points,
        \"notes\": \"$notes\"
    }"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/points/redeem")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to check balance
check_balance_interactive() {
    print_section "Check Point Balance"
    
    read -p "Member ID: " member_id
    
    response=$(curl -s "$API_URL/points/balance?member_id=$member_id")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Function to run automated tests
run_automated_tests() {
    print_section "Running Automated Tests"
    
    echo "This will create test data in the database."
    read -p "Continue? (y/n): " confirm
    
    if [ "$confirm" != "y" ]; then
        echo "Cancelled"
        return
    fi
    
    # Create test members
    echo -e "\n${YELLOW}Creating test members...${NC}"
    member1=$(create_test_member "John Doe" "W1TEST" "john@example.com")
    member2=$(create_test_member "Jane Smith" "KD2ABC" "jane@example.com")
    
    member1_id=$(echo "$member1" | jq -r '.id')
    member2_id=$(echo "$member2" | jq -r '.id')
    
    echo "Created member 1 (ID: $member1_id)"
    echo "Created member 2 (ID: $member2_id)"
    
    # Create test events
    echo -e "\n${YELLOW}Creating test events...${NC}"
    event1=$(create_test_event "Winter Field Day" "Annual winter event" 10)
    event2=$(create_test_event "Summer Hamfest" "Summer swap meet" 5)
    
    event1_id=$(echo "$event1" | jq -r '.id')
    event2_id=$(echo "$event2" | jq -r '.id')
    
    echo "Created event 1 (ID: $event1_id)"
    echo "Created event 2 (ID: $event2_id)"
    
    # Allocate points
    echo -e "\n${YELLOW}Allocating points...${NC}"
    allocate_points "$event1_id" "$member1_id" 15 "Helped with setup"
    allocate_points "$event2_id" "$member1_id" "" "Used default points"
    allocate_points "$event1_id" "$member2_id" 10 "Registration desk"
    
    # Check balances
    echo -e "\n${YELLOW}Checking balances...${NC}"
    balance1=$(check_balance "$member1_id")
    balance2=$(check_balance "$member2_id")
    
    echo "Member 1 balance: $balance1"
    echo "Member 2 balance: $balance2"
    
    echo -e "\n${GREEN}Automated tests completed!${NC}"
}

# Main loop
main() {
    # Check if server is running
    echo -n "Checking server connection... "
    if curl -s "$API_URL/members" > /dev/null 2>&1; then
        echo -e "${GREEN}Connected${NC}"
    else
        echo -e "${RED}Failed${NC}"
        echo "Server is not running at $BASE_URL"
        echo "Start the server with: ./rara-membership serve"
        exit 1
    fi
    
    while true; do
        show_menu
        read -p "Enter your choice: " choice
        
        case $choice in
            1) list_members ;;
            2) create_member_interactive ;;
            3) update_member_interactive ;;
            4) get_member_points ;;
            5) list_events ;;
            6) create_event_interactive ;;
            7) update_event_interactive ;;
            8) allocate_points_interactive ;;
            9) redeem_points_interactive ;;
            10) check_balance_interactive ;;
            11) run_automated_tests ;;
            0) echo "Goodbye!"; exit 0 ;;
            *) echo -e "${RED}Invalid option${NC}" ;;
        esac
        
        echo ""
        read -p "Press Enter to continue..."
    done
}

# Run main function
main