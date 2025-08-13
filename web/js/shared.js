// Shared Alpine.js Components and Functions

// Global function for viewPoints - accessible from Alpine buttons
window.globalViewPoints = function(callSign) {
    window.location.href = `points.html?member=${encodeURIComponent(callSign)}`;
};

// Member Modal Component
function memberModal() {
    return {
        isOpen: false,
        isEditing: false,
        submitting: false,
        error: '',
        success: '',
        originalCallSign: '', // Store original call sign for API calls
        formData: {
            name: '',
            call_sign: '',
            email: '',
            phone: '',
            mailing_address: '',
            frn: '',
            license_class: '',
            member_type: '',
            is_active: true,
            is_silent_key: false
        },
        
        init() {
            // Register this instance globally
            window.memberModalInstance = this;
        },
        
        showCreateModal() {
            this.isEditing = false;
            this.resetForm();
            this.isOpen = true;
        },
        
        showEditModal(member) {
            this.isEditing = true;
            this.originalCallSign = member.call_sign || '';
            this.formData = {
                name: member.name || '',
                call_sign: member.call_sign || '',
                email: member.email || '',
                phone: member.phone || '',
                mailing_address: member.mailing_address || '',
                frn: member.frn || '',
                license_class: member.license_class || '',
                member_type: member.member_type || '',
                is_active: member.is_active !== false,
                is_silent_key: member.is_silent_key === true
            };
            this.isOpen = true;
        },
        
        close() {
            this.isOpen = false;
            this.error = '';
            this.success = '';
        },
        
        resetForm() {
            this.formData = {
                name: '',
                call_sign: '',
                email: '',
                phone: '',
                mailing_address: '',
                frn: '',
                license_class: '',
                member_type: '',
                is_active: true,
                is_silent_key: false
            };
        },
        
        async submitForm() {
            this.submitting = true;
            this.error = '';
            this.success = '';
            
            try {
                const url = this.isEditing 
                    ? `/api/members/${this.originalCallSign}`
                    : '/api/members';
                
                const method = this.isEditing ? 'PUT' : 'POST';
                
                // Prepare data for submission
                const submitData = { ...this.formData };
                
                // Clean up empty strings
                Object.keys(submitData).forEach(key => {
                    if (submitData[key] === '') {
                        submitData[key] = null;
                    }
                });
                
                const response = await fetch(url, {
                    method: method,
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(submitData)
                });
                
                if (!response.ok) {
                    const errorText = await response.text();
                    console.error('API Error:', response.status, errorText);
                    console.error('Submitted data:', JSON.stringify(submitData, null, 2));
                    throw new Error(errorText || `Failed to save member (${response.status})`);
                }
                
                this.success = this.isEditing 
                    ? 'Member updated successfully!' 
                    : 'Member created successfully!';
                
                // Refresh the current page
                setTimeout(() => {
                    window.location.reload();
                }, 1500);
                
            } catch (error) {
                console.error('Submit error:', error);
                this.error = error.message || 'An error occurred while saving the member.';
            } finally {
                this.submitting = false;
            }
        }
    }
}

// Event Modal Component
function eventModal() {
    return {
        isOpen: false,
        isEditing: false,
        submitting: false,
        error: '',
        success: '',
        formData: {
            name: '',
            description: '',
            date: '',
            timezone: 'America/New_York',
            points: 10,
            latitude: 43.138,
            longitude: -77.572
        },
        
        init() {
            // Register this instance globally
            window.eventModalInstance = this;
        },
        
        showCreateModal() {
            this.isEditing = false;
            this.resetForm();
            this.isOpen = true;
        },
        
        showEditModal(event) {
            this.isEditing = true;
            this.formData = {
                id: event.id,
                name: event.name || '',
                description: event.description || '',
                date: event.date ? new Date(event.date).toISOString().slice(0, 16) : '',
                timezone: event.timezone || 'America/New_York',
                points: event.default_points_allocated || 0,
                latitude: event.latitude || 0,
                longitude: event.longitude || 0
            };
            this.isOpen = true;
        },
        
        close() {
            this.isOpen = false;
            this.error = '';
            this.success = '';
        },
        
        resetForm() {
            this.formData = {
                name: '',
                description: '',
                date: '',
                timezone: 'America/New_York',
                points: 10,
                latitude: 43.138,
                longitude: -77.572
            };
        },
        
        async submitForm() {
            this.submitting = true;
            this.error = '';
            this.success = '';
            
            try {
                const url = this.isEditing 
                    ? `/api/events/${this.formData.id}`
                    : '/api/events';
                
                const method = this.isEditing ? 'PUT' : 'POST';
                
                // Prepare data for submission
                const submitData = { ...this.formData };
                delete submitData.id; // Remove id from submission data
                
                // Convert datetime-local to UTC considering the selected timezone
                if (submitData.date && submitData.timezone) {
                    try {
                        // Ensure we have a full datetime string
                        const dateTimeString = submitData.date.includes('T') 
                            ? submitData.date 
                            : submitData.date + 'T12:00';
                        
                        // Add seconds if missing
                        const fullDateTime = dateTimeString.includes(':') && dateTimeString.split(':').length === 2
                            ? dateTimeString + ':00'
                            : dateTimeString;
                        
                        // Create a date as if it's in UTC first
                        const utcDate = new Date(fullDateTime + 'Z');
                        
                        // Format this date in the target timezone to get the actual local time
                        const targetFormatter = new Intl.DateTimeFormat('sv-SE', { // Use Swedish format for ISO-like output
                            timeZone: submitData.timezone,
                            year: 'numeric',
                            month: '2-digit',
                            day: '2-digit',
                            hour: '2-digit',
                            minute: '2-digit',
                            second: '2-digit',
                            hour12: false
                        });
                        
                        const targetTimeString = targetFormatter.format(utcDate);
                        const targetDate = new Date(targetTimeString);
                        
                        // Calculate the offset between what we want and what we got
                        const inputDateTime = new Date(fullDateTime);
                        const offsetMillis = inputDateTime.getTime() - targetDate.getTime();
                        
                        // Apply the offset to get the correct UTC time
                        const correctUtcDate = new Date(utcDate.getTime() + offsetMillis);
                        submitData.date = correctUtcDate.toISOString().slice(0, 16);
                        
                        console.log('Timezone conversion:', {
                            input: fullDateTime,
                            timezone: submitData.timezone,
                            result: submitData.date + 'Z'
                        });
                    } catch (error) {
                        console.error('Error converting datetime with timezone:', error);
                        // Fallback: keep original date
                    }
                }
                
                // Clean up empty strings first
                Object.keys(submitData).forEach(key => {
                    if (submitData[key] === '') {
                        submitData[key] = null;
                    }
                });
                
                // Convert numeric fields to numbers (after cleaning empty strings)
                if (submitData.points !== null && submitData.points !== undefined && submitData.points !== '') {
                    submitData.points = parseInt(submitData.points, 10);
                }
                if (submitData.latitude !== null && submitData.latitude !== undefined && submitData.latitude !== '') {
                    submitData.latitude = parseFloat(submitData.latitude);
                }
                if (submitData.longitude !== null && submitData.longitude !== undefined && submitData.longitude !== '') {
                    submitData.longitude = parseFloat(submitData.longitude);
                }
                
                const response = await fetch(url, {
                    method: method,
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(submitData)
                });
                
                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(errorText || 'Failed to save event');
                }
                
                this.success = this.isEditing 
                    ? 'Event updated successfully!' 
                    : 'Event created successfully!';
                
                // Refresh the current page
                setTimeout(() => {
                    window.location.reload();
                }, 1500);
                
            } catch (error) {
                console.error('Submit error:', error);
                this.error = error.message || 'An error occurred while saving the event.';
            } finally {
                this.submitting = false;
            }
        }
    }
}

// Points Allocation Modal Component
function pointsAllocationModal() {
    return {
        isOpen: false,
        submitting: false,
        error: '',
        success: '',
        memberCallSign: '',
        events: [],
        formData: {
            event_id: '',
            points: null,
            notes: ''
        },
        
        init() {
            // Register this instance globally
            window.pointsAllocationModalInstance = this;
            // Load events when component initializes
            this.loadEvents();
        },
        
        async loadEvents() {
            try {
                const response = await fetch('/api/events');
                this.events = await response.json();
            } catch (error) {
                console.error('Failed to load events:', error);
                this.events = [];
            }
        },
        
        showModal(callSign) {
            this.memberCallSign = callSign;
            this.resetForm();
            this.isOpen = true;
        },
        
        close() {
            this.isOpen = false;
            this.error = '';
            this.success = '';
        },
        
        resetForm() {
            this.formData = {
                event_id: '',
                points: null,
                notes: ''
            };
        },
        
        async submitForm() {
            this.submitting = true;
            this.error = '';
            this.success = '';
            
            try {
                const submitData = {
                    event_id: this.formData.event_id,
                    notes: this.formData.notes || ''
                };
                
                // Only include points if specified
                if (this.formData.points !== null && this.formData.points !== '') {
                    submitData.points = parseInt(this.formData.points, 10);
                }
                
                const response = await fetch(`/api/members/${this.memberCallSign}/points`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(submitData)
                });
                
                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(errorText || 'Failed to allocate points');
                }
                
                this.success = 'Points allocated successfully!';
                
                // Refresh the current page, preserving member selection
                setTimeout(() => {
                    if (window.location.pathname.includes('points.html')) {
                        // Stay on points page with member preserved in URL hash
                        window.location.hash = `member=${encodeURIComponent(this.memberCallSign)}`;
                        window.location.reload();
                    } else {
                        // Redirect to points page with member selection
                        window.location.href = `points.html#member=${encodeURIComponent(this.memberCallSign)}`;
                    }
                }, 1500);
                
            } catch (error) {
                console.error('Submit error:', error);
                this.error = error.message || 'An error occurred while allocating points.';
            } finally {
                this.submitting = false;
            }
        }
    }
}

// Points Redemption Modal Component
function pointsRedemptionModal() {
    return {
        isOpen: false,
        submitting: false,
        error: '',
        success: '',
        memberCallSign: '',
        currentBalance: 0,
        formData: {
            points: 1,
            notes: ''
        },
        
        init() {
            // Register this instance globally
            window.pointsRedemptionModalInstance = this;
        },
        
        showModal(callSign, balance) {
            this.memberCallSign = callSign;
            this.currentBalance = balance;
            this.resetForm();
            this.isOpen = true;
        },
        
        close() {
            this.isOpen = false;
            this.error = '';
            this.success = '';
        },
        
        resetForm() {
            this.formData = {
                points: 1,
                notes: ''
            };
        },
        
        async submitForm() {
            if (this.formData.points > this.currentBalance) {
                this.error = 'Cannot redeem more points than available balance';
                return;
            }
            
            this.submitting = true;
            this.error = '';
            this.success = '';
            
            try {
                const response = await fetch(`/api/members/${this.memberCallSign}/redemptions`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        points: parseInt(this.formData.points, 10),
                        notes: this.formData.notes
                    })
                });
                
                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(errorText || 'Failed to redeem points');
                }
                
                this.success = 'Points redeemed successfully!';
                
                // Refresh the current page, preserving member selection
                setTimeout(() => {
                    if (window.location.pathname.includes('points.html')) {
                        // Stay on points page with member preserved in URL hash
                        window.location.hash = `member=${encodeURIComponent(this.memberCallSign)}`;
                        window.location.reload();
                    } else {
                        // Redirect to points page with member selection
                        window.location.href = `points.html#member=${encodeURIComponent(this.memberCallSign)}`;
                    }
                }, 1500);
                
            } catch (error) {
                console.error('Submit error:', error);
                this.error = error.message || 'An error occurred while redeeming points.';
            } finally {
                this.submitting = false;
            }
        }
    }
}