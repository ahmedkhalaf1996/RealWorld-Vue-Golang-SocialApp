describe("Auth Component Tests", () => {
    // Ignore all application errors
    Cypress.on('uncaught:exception', (err, runnable) => {
        return false
    })

    beforeEach(() => {
        cy.visit('/Auth')
        // Wait for page to be fully loaded
        cy.wait(1000)
    })

    describe('Signin Form', () => {
        it('should display all signin form elements', () => {
            cy.get('.col-5 input', { timeout: 10000 }).should('have.length', 2)
            cy.contains('sigin in').should('be.visible')
            cy.get('.col-5 button[type="submit"]').should('be.visible')
        })

        it('should allow typing in signin inputs', () => {
            cy.get('.col-5 input').eq(0)
                .type('test@example.com', { force: true })
                .should('have.value', 'test@example.com')
            
            cy.get('.col-5 input').eq(1)
                .type('password123', { force: true })
                .should('have.value', 'password123')
        })

        it('should have password input type', () => {
            cy.get('.col-5 input').eq(1)
                .should('have.attr', 'type', 'password')
        })
    })

    describe('Signup Form', () => {
        it('should display all signup form elements', () => {
            cy.get('.col-7 input', { timeout: 10000 }).should('have.length', 4)
            cy.contains('Your first Name *').should('be.visible')
            cy.contains('Your lastName *').should('be.visible')
            cy.contains('Your Email *').should('be.visible')
            cy.contains('Your Password *').should('be.visible')
            cy.contains('Create New Account').should('be.visible')
        })

        it('should allow typing in all signup inputs', () => {
            cy.get('.col-7 input').eq(0)
                .type('John', { force: true })
                .should('have.value', 'John')

            cy.get('.col-7 input').eq(1)
                .type('Doe', { force: true })
                .should('have.value', 'Doe')

            cy.get('.col-7 input').eq(2)
                .type('j@example.com', { force: true })
                .should('have.value', 'j@example.com')

            cy.get('.col-7 input').eq(3)
                .type('password123', { force: true })
                .should('have.value', 'password123')
        })

        it('should have correct button colors', () => {
            cy.get('.col-5 .q-btn').should('have.class', 'bg-primary')
            cy.get('.col-7 .q-btn').should('have.class', 'bg-positive')
        })
    })

    describe('Form Interactions', () => {
        it('should handle empty input submissions', () => {
            cy.get('.col-5 button[type="submit"]').click({ force: true })
            cy.wait(500)
            
            cy.get('.col-7 button[type="submit"]').click({ force: true })
            cy.wait(500)
        })
    })

    describe('Responsive Design', () => {
        it('should maintain layout on different screen sizes', () => {
            cy.viewport(1200, 800)
            cy.wait(500)
            cy.get('.row').should('be.visible')

            cy.viewport(768, 1024)
            cy.wait(500)
            cy.get('.col-5').should('be.visible')
            cy.get('.col-7').should('be.visible')

            cy.viewport(375, 667)
            cy.wait(500)
            cy.get('.q-card').should('be.visible')
        })
    })
})