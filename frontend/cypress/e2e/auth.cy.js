describe("Auth Component Tests", () => {
    beforeEach(() => {
        cy.visit('/Auth')
    })


    describe('Page Layout', () => {
        it('should switch between signin and signup tabs on mobile', () => {
            cy.viewport(375, 667)
            cy.get('.mobile-tabs').should('be.visible')

            // Click signup tab
            cy.contains('Sign Up').click()
            cy.get('[name="signup"]').should('be.visible')

            // Click signin tab
            cy.contains('Sign In').click()
            cy.get('[name="signin"]').should('be.visible')
        })
    })

    describe('Signin Form', () => {
        it('should display all signin form elements', () => {
            cy.get('.col-5 input').should('have.length', 2)

            // chek for buttom
            cy.contains('sigin in').should('be.visible')
            cy.get('.col-5 button[type="submit"]').should('be.visible')
        })

        it('should allow typing in signin inputs', () => {
            cy.get('.col-5 input').eq(0)
                .type('test@example.com')
                .should('have.value', 'test@example.com')

            cy.get('.col-5 input').eq(1)
                .type('password123')
                .should('have.value', 'password123')
        })

        it('should have password input type', () => {
            cy.get('.col-5 input').eq(1)
                .should('have.attr', 'type', 'password')
        })
    })

    describe('Signup Form', () => {
        it('should display all signup form elemnts', () => {
            // chekc all s inputs exists
            cy.get('.col-7 input').should('have.length', 4)
            cy.contains('Your first Name *').should('be.visible')
            cy.contains('Your lastName *').should('be.visible')
            cy.contains('Your Email *').should('be.visible')
            cy.contains('Your Password *').should('be.visible')

            // check s button
            cy.contains('Create New Account').should('be.visible')
        })

        it('should allow typing in all signup inputs', () => {
            cy.get('.col-7 input').eq(0)
                .type('John')
                .should('have.value', 'John')

            cy.get('.col-7 input').eq(1)
                .type('Doe')
                .should('have.value', 'Doe')

            cy.get('.col-7 input').eq(2)
                .type('j@example.com')
                .should('have.value', 'j@example.com')

            cy.get('.col-7 input').eq(3)
                .type('passowrd123')
                .should('have.value', 'passowrd123')
        })

        it('should have correct buttom colors', () => {
            // signin
            cy.get('.col-5 .q-btn').should('have.class', 'bg-primary')

            // signup
            cy.get('.col-7 .q-btn').should('have.class', 'bg-positive')

        })
    })

    describe('Form Interactions', () => {
        it('Should Handle Empty input Summintions', () => {
            cy.get('.col-5 button[type="submit"]').click()
            cy.get('.q-notification').should('be.visible').and('contain', 'Email is Required')

            cy.get('.col-7 button[type="submit"]').click()
            cy.get('.q-notification').should('be.visible').and('contain', 'Email is Required')
            cy.get('.q-notification').should('be.visible').and('contain', 'password is Required')
            cy.get('.q-notification').should('be.visible').and('contain', 'firstName is Required')
            cy.get('.q-notification').should('be.visible').and('contain', 'lastName is Required')
        })
    })

    describe('Responsive Desing', () => {
        it('should maintian layout on diffrent screen sizes', () => {
            cy.viewport(375, 667)
            cy.get('.q-card').should('be.visible')

            cy.viewport(768, 1024)
            cy.get('.col-5').should('be.visible')
            cy.get('.col-7').should('be.visible')

            cy.viewport(1200, 800)
            cy.get('.row').should('be.visible')
        })
    })
})
