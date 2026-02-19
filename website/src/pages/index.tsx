import type {ReactNode} from 'react';
import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';
import CodeBlock from '@theme/CodeBlock';

function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    const split = siteConfig.tagline.split(" ")

    return (
        <header className={clsx('hero hero--primary', styles.heroBanner)}>
            <div className="container">
                <Heading as="h1" className="hero__title">
                    {siteConfig.title}
                </Heading>
                <p
                    className="hero__subtitle"
                    style={{color: "white"}}
                >
                  <span className="text--primary">
                    {split[0]}
                  </span>{" "}
                    {split.slice(1).join(" ")}
                </p>

                <TerminalCommand command="docker run --rm -p 6699:6699 ghcr.io/ra341/glacier:canary"/>

                <div className={styles.buttons}>
                    <Link
                        className="button button--secondary button--lg"
                        to="/docs/intro">
                        Friendly Manual
                    </Link>
                </div>
            </div>
        </header>
    );
}

export default function Home(): ReactNode {
    const {siteConfig} = useDocusaurusContext();
    return (
        <Layout
            title={`Hello from ${siteConfig.title}`}
            description="Description will go into a meta tag in <head />">
            <HomepageHeader/>
            <main>
                <HomepageFeatures/>
            </main>
        </Layout>
    );
}

function TerminalCommand({command}) {
    return (
        <div className="terminal-container" style={{
            maxWidth: '600px',
            margin: '2rem auto 0',
            textAlign: 'left',
            borderRadius: '8px',
            marginBottom: '2rem',
            overflow: 'hidden',
            boxShadow: '0 10px 30px rgba(0,0,0,0.3)',
            backgroundColor: '#1c1e21'
        }}>
            {/* Terminal Header Bar */}
            <div style={{
                backgroundColor: '#303846',
                padding: '10px 15px',
                display: 'flex',
                alignItems: 'center',
                gap: '8px'
            }}>
                <span style={{
                    marginLeft: '10px',
                    color: '#a0a0a0',
                    fontSize: '0.8rem',
                    fontFamily: 'monospace'
                }}>
                    Try now
                </span>
            </div>

            {/* The Command itself */}
            <CodeBlock language="bash">
                {command}
            </CodeBlock>
        </div>
    );
}